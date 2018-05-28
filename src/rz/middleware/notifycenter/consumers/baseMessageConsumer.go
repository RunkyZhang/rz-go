package consumers

import (
	"time"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/global"
	"fmt"
)

type consumeFunc func(jsonString string) (interface{}, error)
type handleErrorFunc func(interface{}, error) (error)

type baseMessageConsumer struct {
	SendChannel enumerations.SendChannel
	consumeFunc consumeFunc
	handleErrorFunc handleErrorFunc
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
	keyValues, err := baseMessageConsumer.getMessages()
	if nil != err {
		fmt.Println("failed to get messages. error:", err)
		return
	}
	if nil == keyValues {
		return
	}

	for key, value := range keyValues {
		messageDto, err := baseMessageConsumer.consumeFunc(value)
		if nil != err {
			fmt.Println("failed to consume message["+key+"]. error:", err)
			if nil != messageDto {
				err = baseMessageConsumer.handleErrorFunc(messageDto, err)
				if nil != err {
					fmt.Println("failed to handle error for message["+key+"]. error:", err)
				}
			}
		}
	}
}

func (baseMessageConsumer *baseMessageConsumer) getMessages() (map[string]string, error) {
	sendChannel, err := enumerations.SendChannelToString(baseMessageConsumer.SendChannel)
	if nil != err {
		return nil, err
	}

	messageIds, err := global.GetRedisClient().SortedSetRangeByScore(global.RedisKeyMessageKeys+sendChannel, 0, float64(time.Now().Unix()))
	if nil != err {
		return nil, err
	}

	return global.GetRedisClient().HashGetMany(global.RedisKeyMessageValues+sendChannel, messageIds...)
}
