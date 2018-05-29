package consumers

import (
	"time"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/global"
	"fmt"
	"rz/middleware/notifycenter/exceptions"
)

type consumeFunc func(string) (interface{}, error)
type handleErrorFunc func(interface{}, error) (error)

type baseMessageConsumer struct {
	SendChannel     enumerations.SendChannel
	consumeFunc     consumeFunc
	handleErrorFunc handleErrorFunc
}

func (baseMessageConsumer *baseMessageConsumer) Start(duration time.Duration) {
	timer := time.NewTimer(duration)

	for {
		select {
		case <-timer.C:
			baseMessageConsumer.start()
			timer.Reset(duration)

			fmt.Println(time.Now())
		}
	}
}

func (baseMessageConsumer *baseMessageConsumer) start() {
	messageIds, err := baseMessageConsumer.getMessageIds()
	if nil != err {
		fmt.Println("failed to get message ids. error:", err)
		return
	}
	if nil == messageIds {
		return
	}

	//return

	for _, messageId := range messageIds {
		var messageDto interface{}
		jsonString, err := global.GetRedisClient().HashGet(global.RedisKeyMessageValues+sendChannel, messageId)
		if nil != err {
			fmt.Println("failed to get message by id[", messageId, "]. error:", err)
			err = exceptions.MessageBodyMissed
		} else {
			messageDto, err = baseMessageConsumer.consumeFunc(jsonString)
		}

		if nil != err {
			fmt.Println("failed to consume message. error:", err)
			if nil != messageDto {
				err = baseMessageConsumer.handleErrorFunc(messageDto, err)
				if nil != err {
					fmt.Println("failed to handle error for message. error:", err)
				}
			}
		}
	}
}

func (baseMessageConsumer *baseMessageConsumer) getMessageIds() ([]string, error) {
	sendChannel, err := enumerations.SendChannelToString(baseMessageConsumer.SendChannel)
	if nil != err {
		return nil, err
	}

	max := float64(time.Now().Unix())
	messageIds, err := global.GetRedisClient().SortedSetRangeByScore(global.RedisKeyMessageKeys+sendChannel, 0, max)
	if nil != err {
		return nil, err
	}

	return messageIds, nil
}

func ConsumerStart() {
	duration := time.Duration(global.Config.ConsumingInterval) * time.Second
	for i := 0; i < 1; i++ {
		//go MailMessageConsumer.Start(duration)
		go SmsMessageConsumer.Start(duration)
		time.Sleep(2 * time.Second)
	}
}
