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
)

type convertFunc func(string) (interface{}, *models.BaseMessageDto, error)
type sendFunc func(interface{}) (error)

type baseMessageConsumer struct {
	SendChannel enumerations.SendChannel
	keySuffix   string
	convertFunc convertFunc
	sendFunc    sendFunc
}

func (baseMessageConsumer *baseMessageConsumer) Start(duration time.Duration) {
	timer := time.NewTimer(duration)

	for {
		select {
		case <-timer.C:
			baseMessageConsumer.start()
			timer.Reset(duration)
		}
	}
}

func (baseMessageConsumer *baseMessageConsumer) start() {
	messageIds, err := baseMessageConsumer.getMessageIds()
	if nil != err || nil == messageIds {
		fmt.Println("failed to get message ids. error: ", err)
		return
	}

	for _, messageId := range messageIds {
		jsonString, err := global.GetRedisClient().HashGet(global.RedisKeyMessageValues+baseMessageConsumer.keySuffix, messageId)

		if nil == err {
			var messageDto interface{}
			var baseMessageDto *models.BaseMessageDto
			var flagError error

			count, err := global.GetRedisClient().SortedSetRemoveByValue(global.RedisKeyMessageKeys+baseMessageConsumer.keySuffix, messageId)
			if nil == err {
				if 0 < count {
					messageDto, baseMessageDto, flagError = baseMessageConsumer.convertFunc(jsonString)
					if nil == flagError {
						flagError = baseMessageConsumer.consume(messageDto, baseMessageDto)
						if nil == flagError {
							fmt.Println("success to consume message[", messageId, "]")
						}
					}
				}
			} else {
				fmt.Println("failed to remove message[", messageId, "]. error: ", err)
			}

			if nil != flagError {
				fmt.Println("failed to consume message[", messageId, "]. error: ", flagError)

				// when string is error json string
				if nil != messageDto {
					err = baseMessageConsumer.handleError(messageDto, baseMessageDto, flagError)
					if nil != err {
						fmt.Println("failed to handle error for message[", messageId, "]. error: ", err)
					}
				}
			}
		} else {
			// ignore message
			fmt.Println("failed to get message[", messageId, "] value. error: ", err)

			_, err := global.GetRedisClient().SortedSetRemoveByValue(global.RedisKeyMessageKeys+baseMessageConsumer.keySuffix, messageId)
			if nil != err {
				fmt.Println("failed to remove message[", messageId, "]. error: ", err)
			}
		}
	}
}

func (baseMessageConsumer *baseMessageConsumer) getMessageIds() ([]string, error) {
	max := float64(time.Now().Unix())
	messageIds, err := global.GetRedisClient().SortedSetRangeByScore(global.RedisKeyMessageKeys+baseMessageConsumer.keySuffix, 0, max)
	if nil != err {
		return nil, err
	}

	return messageIds, nil
}

func (baseMessageConsumer *baseMessageConsumer) consume(messageDto interface{}, baseMessageDto *models.BaseMessageDto) (error) {
	baseMessageDto.States = baseMessageDto.States + "+" + enumerations.MessageStateToString(enumerations.Consuming)
	baseMessageConsumer.updateMessage(messageDto, baseMessageDto.Id)

	if time.Now().Unix() > baseMessageDto.ExpireTime {
		return exceptions.MessageExpire
	}

	err := baseMessageConsumer.sendFunc(messageDto)
	if nil == err {
		baseMessageDto.States = baseMessageDto.States + "+" + enumerations.MessageStateToString(enumerations.Sent)
		baseMessageDto.Finished = true
		baseMessageConsumer.updateMessage(messageDto, baseMessageDto.Id)
	}

	return err
}

func (baseMessageConsumer *baseMessageConsumer) handleError(messageDto interface{}, baseMessageDto *models.BaseMessageDto, flagError error) (error) {
	var messageState string
	if flagError == exceptions.MessageExpire {
		messageState = enumerations.MessageStateToString(enumerations.Expire)
	} else {
		messageState = enumerations.MessageStateToString(enumerations.Error)
	}

	baseMessageDto.States = baseMessageDto.States + "+" + messageState
	baseMessageDto.ErrorMessage = flagError.Error()
	baseMessageDto.Finished = true

	return baseMessageConsumer.updateMessage(messageDto, baseMessageDto.Id)
}

func (baseMessageConsumer *baseMessageConsumer) updateMessage(messageDto interface{}, messageId string) (error) {
	bytes, err := json.Marshal(messageDto)
	if nil != err {
		return err
	}

	return global.GetRedisClient().HashSet(global.RedisKeyMessageValues+baseMessageConsumer.keySuffix, messageId, string(bytes))
}

func ConsumerStart() {
	duration := time.Duration(global.Config.ConsumingInterval) * time.Second
	for i := 0; i < runtime.NumCPU(); i++ {
		go MailMessageConsumer.Start(duration)
		go SmsMessageConsumer.Start(duration)
		time.Sleep(2 * time.Second)
	}
}
