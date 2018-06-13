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
)

type convertFunc func(int, time.Time) (interface{}, *models.MessageBasePo, error)
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
	if nil != err || nil == messageIds {
		fmt.Println("failed to get message ids. error: ", err)
		return
	}

	for _, messageId := range messageIds {
		if nil != err {
			// ignore message
			fmt.Printf("failed to get message(%d) value. error: %s", messageId, err.Error())

			affectedCount, err := myself.messageManagementBase.RemoveMessageId(messageId)
			if nil != err || 0 == affectedCount {
				fmt.Println("failed to remove message(", messageId, ") value. error: ", err)
			}

			continue
		}

		affectedCount, err := myself.messageManagementBase.RemoveMessageId(messageId)
		if nil != err || 0 == affectedCount {
			fmt.Println("failed to remove message(", messageId, "). error: ", err)
			continue
		}

		var messagePo interface{}
		var messageBasePo *models.MessageBasePo
		var flagError error
		messagePo, messageBasePo, flagError = myself.convertFunc(messageId, now)
		if nil == flagError {
			flagError = myself.consume(messagePo, messageBasePo)
			if nil == flagError {
				fmt.Printf("success to consume message(%d)", messageId)
			}
		}

		if nil != flagError {
			fmt.Printf("failed to consume message(%d). error: %s", messageId, flagError.Error())

			// when string is error json string
			if nil != messageBasePo {
				var state string
				if flagError == exceptions.MessageExpire {
					state = enumerations.MessageStateToString(enumerations.Expire)
				} else {
					state = enumerations.MessageStateToString(enumerations.Error)
				}
				myself.modifyMessagePo(messageBasePo, state, true, flagError.Error())
			}
		}
	}
}

func (myself *messageConsumerBase) consume(messagePo interface{}, messageBasePo *models.MessageBasePo) (error) {
	myself.modifyMessagePo(messageBasePo, enumerations.MessageStateToString(enumerations.Consuming), false, "")

	if time.Now().Unix() > messageBasePo.ExpireTime.Unix() {
		return exceptions.MessageExpire
	}

	err := myself.sendFunc(messagePo)
	if nil != err {
		return err
	}

	myself.modifyMessagePo(messageBasePo, enumerations.MessageStateToString(enumerations.Sent), true, "")
	return nil
}

func (myself *messageConsumerBase) modifyMessagePo(messageBasePo *models.MessageBasePo, state string, finished bool, errorMessage string) {
	messageBasePo.States = messageBasePo.States + "+" + state
	var errorMessages string
	if "" == errorMessage {
		errorMessages = ""
	} else {
		messageBasePo.ErrorMessages = messageBasePo.ErrorMessages + "+++" + errorMessage
		errorMessages = messageBasePo.ErrorMessages
	}

	affectedCount, err := myself.messageManagementBase.ModifyById(
		messageBasePo.Id,
		messageBasePo.States,
		finished,
		errorMessages,
		messageBasePo.CreatedTime)
	if nil != err || 0 == affectedCount {
		fmt.Println("failed to modfiy message(", messageBasePo.Id, ") state. error: ", err.Error())
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
