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

	key := global.RedisKeyMessage + sendChannel
	score := messageDto.CreatedTime

	err = global.GetRedisClient().HashSet(key+"_values", messageDto.Id, jsonString)
	if nil != err {
		return err
	}

	return global.GetRedisClient().SortedSetAdd(key+"_keys", messageDto.Id, float64(score))
}
