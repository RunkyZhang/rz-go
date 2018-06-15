package consumers

import (
	"time"
	"fmt"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/common"
	"runtime"
	"errors"
)

type convertFunc func(int, time.Time) (interface{}, *models.PoBase, *models.CallbackBasePo, error)
type sendFunc func(interface{}) (error)

type messageConsumerBase struct {
	convertFunc           convertFunc
	sendFunc              sendFunc
	messageManagementBase *managements.MessageManagementBase
	asyncJobWorker        *common.AsyncJobWorker
}

func (myself *messageConsumerBase) Start(duration time.Duration) (error) {
	if nil == myself.messageManagementBase {
		return errors.New("[myself.messageManagementBase] is nil")
	}

	asyncJob := &common.AsyncJob{
		Name:    myself.messageManagementBase.KeySuffix,
		Type:    "Consumer",
		RunFunc: myself.start,
	}
	myself.asyncJobWorker = common.NewAsyncJobWorker(runtime.NumCPU(), duration, asyncJob)

	return nil
}

func (myself *messageConsumerBase) start(parameter interface{}) (error) {
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

			managements.ModifyMessageFlowAsync(
				myself.messageManagementBase,
				messageId,
				poBase,
				callbackBasePo,
				messageState,
				true,
				flagError.Error())
		}
	}

	return nil
}

func (myself *messageConsumerBase) consume(messagePo interface{}, messageId int, poBase *models.PoBase, callbackBasePo *models.CallbackBasePo) (error) {
	managements.ModifyMessageFlowAsync(
		myself.messageManagementBase,
		messageId,
		poBase,
		callbackBasePo,
		enumerations.Consuming,
		false,
		"")

	if time.Now().Unix() > callbackBasePo.ExpireTime.Unix() {
		return exceptions.MessageExpire().AttachMessage(common.Int32ToString(messageId))
	}

	err := myself.sendFunc(messagePo)
	if nil != err {
		return err
	}

	managements.ModifyMessageFlowAsync(
		myself.messageManagementBase,
		messageId,
		poBase,
		callbackBasePo,
		enumerations.Sent,
		true,
		"")

	return nil
}

//func ConsumerStart() {
//	duration := time.Duration(global.Config.ConsumingInterval) * time.Second
//	for i := 0; i < runtime.NumCPU(); i++ {
//		go MailMessageConsumer.Start(duration)
//		go SmsMessageConsumer.Start(duration)
//		time.Sleep(2 * time.Second)
//	}
//}
