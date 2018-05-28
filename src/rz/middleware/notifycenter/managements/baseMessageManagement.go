package managements

import (
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
)

type baseMessageManagement struct {
}

func (*baseMessageManagement) addMessage(messageDto *models.MessageDto, jsonString string) (error) {
	sendChannel, err := enumerations.SendChannelToString(messageDto.SendChannel)
	if nil != err {
		return err
	}

	err = global.GetRedisClient().HashSet(global.RedisKeyMessageValues+sendChannel, messageDto.Id, jsonString)
	if nil != err {
		return err
	}

	return global.GetRedisClient().SortedSetAdd(global.RedisKeyMessageKeys+sendChannel, messageDto.Id, float64(messageDto.CreatedTime))
}
