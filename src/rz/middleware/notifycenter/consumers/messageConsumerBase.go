package consumers

import (
	"time"
	"fmt"
	"runtime"

	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/common"
)

type convertFunc func(int, time.Time) (interface{}, *models.PoBase, *models.CallbackBasePo, error)
type sendFunc func(interface{}) (error)

type messageConsumerBase struct {
	convertFunc           convertFunc
	sendFunc              sendFunc
	messageManagementBase *managements.MessageManagementBase
}

func (myself *messageConsumerBase) Start(duration time.Duration) {
	timer := time.NewTimer(duration)

	for {
		select {
		case <-timer.C:
			myself.start()
			timer.Reset(duration)
		}
	}
}

func (myself *messageConsumerBase) start() {
	now := time.Now()
	messageIds, err := myself.messageManagementBase.DequeueMessageIds(now)
	if nil != err {
		fmt.Printf("failed to get message(%s) ids. error: %s", myself.messageManagementBase.KeySuffix, err.Error())
		return
	}
	if nil == messageIds {
		return
	}

	for _, messageId := range messageIds {
		affectedCount, err := myself.messageManagementBase.RemoveMessageId(messageId)
		if nil != err || 0 == affectedCount {
			fmt.Printf("failed to remove message(%d). error: %s", messageId, err)
			continue
		}

		messagePo, poBase, callbackBasePo, flagError := myself.convertFunc(messageId, now)
		if nil == flagError {
			flagError = myself.consume(messagePo, messageId, poBase, callbackBasePo)
			if nil == flagError {
				fmt.Printf("success to consume message(%d)", messageId)
			}
		}

		if nil != flagError {
			fmt.Printf("failed to consume message(%d). error: %s", messageId, flagError.Error())

			var messageState enumerations.MessageState
			businessError, ok := flagError.(*exceptions.BusinessError)
			if ok && exceptions.MessageExpire().Code == businessError.Code {
				messageState = enumerations.Expire
			} else {
				messageState = enumerations.Error
			}

			myself.modifyMessageFlow(messageId, poBase, callbackBasePo, messageState, true, flagError.Error())
		}
	}
}

func (myself *messageConsumerBase) consume(messagePo interface{}, messageId int, poBase *models.PoBase, callbackBasePo *models.CallbackBasePo) (error) {
	myself.modifyMessageFlow(messageId, poBase, callbackBasePo, enumerations.Consuming, false, "")

	if time.Now().Unix() > callbackBasePo.ExpireTime.Unix() {
		return exceptions.MessageExpire().AttachMessage(common.Int32ToString(messageId))
	}

	err := myself.sendFunc(messagePo)
	if nil != err {
		return err
	}

	myself.modifyMessageFlow(messageId, poBase, callbackBasePo, enumerations.Sent, true, "")
	return nil
}

func (myself *messageConsumerBase) modifyMessageFlow(
	messageId int,
	poBase *models.PoBase,
	callbackBasePo *models.CallbackBasePo,
	messageState enumerations.MessageState,
	finished bool,
	errorMessage string) {
	state := enumerations.MessageStateToString(messageState)
	callbackBasePo.States = callbackBasePo.States + "+" + state
	var errorMessages string
	if "" == errorMessage {
		errorMessages = ""
	} else {
		callbackBasePo.ErrorMessages = callbackBasePo.ErrorMessages + "+++" + errorMessage
		errorMessages = callbackBasePo.ErrorMessages
	}

	affectedCount, err := myself.messageManagementBase.ModifyById(
		messageId,
		callbackBasePo.States,
		finished,
		errorMessages,
		poBase.CreatedTime)
	if nil != err || 0 == affectedCount {
		fmt.Printf("failed to modify message(%d) state. error: %s", messageId, err)
	}
}

func ConsumerStart() {
	duration := time.Duration(global.Config.ConsumingInterval) * time.Second
	for i := 0; i < runtime.NumCPU(); i++ {
		go MailMessageConsumer.Start(duration)
		go SmsMessageConsumer.Start(duration)
		time.Sleep(2 * time.Second)
	}
}
