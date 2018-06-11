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

type convertFunc func(int) (interface{}, *models.MessageBasePo, error)
type sendFunc func(interface{}) (error)

type messageConsumerBase struct {
	convertFunc           convertFunc
	sendFunc              sendFunc
	messageManagementBase *managements.MessageManagementBase
}

func (messageConsumerBase *messageConsumerBase) Start(duration time.Duration) {
	timer := time.NewTimer(duration)

	for {
		select {
		case <-timer.C:
			messageConsumerBase.start()
			timer.Reset(duration)
		}
	}
}

func (messageConsumerBase *messageConsumerBase) start() {
	messageIds, err := messageConsumerBase.messageManagementBase.DequeueMessageIds()
	if nil != err || nil == messageIds {
		fmt.Println("failed to get message ids. error: ", err)
		return
	}

	for _, messageId := range messageIds {
		if nil != err {
			// ignore message
			fmt.Printf("failed to get message(%d) value. error: %s", messageId, err.Error())

			affectedCount, err := messageConsumerBase.messageManagementBase.RemoveMessageId(messageId)
			if nil != err || 0 == affectedCount {
				fmt.Println("failed to remove message(", messageId, ") value. error: ", err)
			}

			continue
		}

		affectedCount, err := messageConsumerBase.messageManagementBase.RemoveMessageId(messageId)
		if nil != err || 0 == affectedCount {
			fmt.Println("failed to remove message(", messageId, "). error: ", err)
			continue
		}

		var messagePo interface{}
		var messageBasePo *models.MessageBasePo
		var flagError error
		messagePo, messageBasePo, flagError = messageConsumerBase.convertFunc(messageId)
		if nil == flagError {
			flagError = messageConsumerBase.consume(messagePo, messageBasePo)
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
				messageConsumerBase.modifyMessagePo(messageBasePo, state, true, flagError.Error())
			}
		}
	}
}

func (messageConsumerBase *messageConsumerBase) consume(messagePo interface{}, messageBasePo *models.MessageBasePo) (error) {
	messageConsumerBase.modifyMessagePo(messageBasePo, enumerations.MessageStateToString(enumerations.Consuming), false, "")

	if time.Now().Unix() > messageBasePo.ExpireTime.Unix() {
		return exceptions.MessageExpire
	}

	err := messageConsumerBase.sendFunc(messagePo)
	if nil != err {
		return err
	}

	messageConsumerBase.modifyMessagePo(messageBasePo, enumerations.MessageStateToString(enumerations.Sent), true, "")
	return nil
}

func (messageConsumerBase *messageConsumerBase) modifyMessagePo(messageBasePo *models.MessageBasePo, state string, finished bool, errorMessage string) {
	messageBasePo.States = messageBasePo.States + "+" + state
	var errorMessages string
	if "" == errorMessage {
		errorMessages = ""
	} else {
		messageBasePo.ErrorMessages = messageBasePo.ErrorMessages + "+++" + errorMessage
		errorMessages = messageBasePo.ErrorMessages
	}

	affectedCount, err := managements.SmsMessageManagement.ModifyById(messageBasePo.Id, messageBasePo.States, finished, errorMessages)
	if nil != err || 0 == affectedCount {
		fmt.Println("failed to handle error for message(", messageBasePo.Id, "). error: ", err.Error())
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
