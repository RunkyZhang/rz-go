package consumers

import (
	"time"
	"fmt"
	"encoding/json"
	"runtime"

	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/common"
)

type convertFunc func(string) (interface{}, *models.MessageBaseDto, error)
type sendFunc func(interface{}) (error)

type messageConsumerBase struct {
	SendChannel enumerations.SendChannel
	keySuffix   string
	convertFunc convertFunc
	sendFunc    sendFunc
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
	messageIds, err := messageConsumerBase.getMessageIds()
	if nil != err || nil == messageIds {
		fmt.Println("failed to get message ids. error: ", err)
		return
	}

	for _, messageId := range messageIds {
		jsonString, err := global.GetRedisClient().HashGet(global.RedisKeyMessageValues+messageConsumerBase.keySuffix, messageId)

		if nil != err {
			// ignore message
			fmt.Printf("failed to get message(%s) value. error: %s", messageId, err.Error())

			_, err := global.GetRedisClient().SortedSetRemoveByValue(global.RedisKeyMessageKeys+messageConsumerBase.keySuffix, messageId)
			if nil != err {
				fmt.Printf("failed to remove message(%s) value. error: %s", messageId, err.Error())
			}

			continue
		}

		count, err := global.GetRedisClient().SortedSetRemoveByValue(global.RedisKeyMessageKeys+messageConsumerBase.keySuffix, messageId)
		if nil != err || 0 == count {
			fmt.Println("failed to remove message[", messageId, "]. error: ", err)
			continue
		}

		var messageDto interface{}
		var messageBaseDto *models.MessageBaseDto
		var flagError error
		messageDto, messageBaseDto, flagError = messageConsumerBase.convertFunc(jsonString)
		if nil == flagError {
			flagError = messageConsumerBase.consume(messageDto, messageBaseDto)
			if nil == flagError {
				fmt.Printf("success to consume message(%s)", messageId)
			}
		}

		if nil != flagError {
			fmt.Printf("failed to consume message(%s). error: %s", messageId, flagError.Error())

			// when string is error json string
			if nil != messageDto {
				err = messageConsumerBase.handleError(messageDto, messageBaseDto, flagError)
				if nil != err {
					fmt.Printf("failed to handle error for message(%s). error: %s", messageId, err.Error())
				}
			}
		}
	}
}

func (messageConsumerBase *messageConsumerBase) getMessageIds() ([]string, error) {
	max := float64(time.Now().Unix())
	messageIds, err := global.GetRedisClient().SortedSetRangeByScore(global.RedisKeyMessageKeys+messageConsumerBase.keySuffix, 0, max)
	if nil != err {
		return nil, err
	}

	return messageIds, nil
}

func (messageConsumerBase *messageConsumerBase) consume(messageDto interface{}, messageBaseDto *models.MessageBaseDto) (error) {
	messageBaseDto.States = messageBaseDto.States + "+" + enumerations.MessageStateToString(enumerations.Consuming)
	messageConsumerBase.updateMessage(messageDto, messageBaseDto.Id)

	if time.Now().Unix() > messageBaseDto.ExpireTime {
		return exceptions.MessageExpire
	}

	err := messageConsumerBase.sendFunc(messageDto)
	if nil == err {
		messageBaseDto.States = messageBaseDto.States + "+" + enumerations.MessageStateToString(enumerations.Sent)
		messageBaseDto.Finished = true
		messageConsumerBase.updateMessage(messageDto, messageBaseDto.Id)
	}

	return err
}

func (messageConsumerBase *messageConsumerBase) handleError(messageDto interface{}, messageBaseDto *models.MessageBaseDto, flagError error) (error) {
	var messageState string
	if flagError == exceptions.MessageExpire {
		messageState = enumerations.MessageStateToString(enumerations.Expire)
	} else {
		messageState = enumerations.MessageStateToString(enumerations.Error)
	}

	messageBaseDto.States = messageBaseDto.States + "+" + messageState
	messageBaseDto.ErrorMessage = flagError.Error()
	messageBaseDto.Finished = true

	return messageConsumerBase.updateMessage(messageDto, messageBaseDto.Id)
}

func (messageConsumerBase *messageConsumerBase) updateMessage(messageDto interface{}, messageId int) (error) {
	bytes, err := json.Marshal(messageDto)
	if nil != err {
		return err
	}

	return global.GetRedisClient().HashSet(global.RedisKeyMessageValues+messageConsumerBase.keySuffix, common.Int32ToString(messageId), string(bytes))
}

func ConsumerStart() {
	duration := time.Duration(global.Config.ConsumingInterval) * time.Second
	for i := 0; i < runtime.NumCPU(); i++ {
		go MailMessageConsumer.Start(duration)
		go SmsMessageConsumer.Start(duration)
		time.Sleep(2 * time.Second)
	}
}
