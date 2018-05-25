package management

import (
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
	"encoding/json"
	"time"
	"rz/middleware/notifycenter/enumerations"
)

var (
	SmsManagement = smsManagement{}
)

type smsManagement struct {
}

func (*smsManagement) SendSms(smsMessageDto *models.SmsMessageDto) (error) {
	sendChannel, err := enumerations.SendChannelToString(smsMessageDto.SendChannel)
	if nil == err {
		return err
	}

	bytes, err := json.Marshal(smsMessageDto)
	if nil == err {
		return err
	}

	key := global.RedisKeyMessage + sendChannel
	value := string(bytes)
	score := time.Now().Unix()

	return global.GetRedisClient().SortedSetAdd(key, value, float64(score))
}
