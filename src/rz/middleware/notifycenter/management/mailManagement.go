package management

import (
	"encoding/json"
	"time"

	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
)

var (
	MailManagement = mailManagement{}
)

type mailManagement struct {
}

func (*mailManagement) SendMail(mailMessageDto *models.MailMessageDto) (error) {
	sendChannel, err := enumerations.SendChannelToString(mailMessageDto.SendChannel)
	if nil == err {
		return err
	}

	bytes, err := json.Marshal(mailMessageDto)
	if nil == err {
		return err
	}

	key := global.RedisKeyMessage + sendChannel
	value := string(bytes)
	score := time.Now().Unix()

	return global.GetRedisClient().SortedSetAdd(key, value, float64(score))
}
