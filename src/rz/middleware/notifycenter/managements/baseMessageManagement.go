package managements

import (
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
)

type baseMessageManagement struct {
}

func (*baseMessageManagement) addMessage(baseMessageDto *models.BaseMessageDto, jsonString string) (error) {
	sendChannel, err := enumerations.SendChannelToString(baseMessageDto.SendChannel)
	if nil != err {
		return err
	}

	err = global.GetRedisClient().HashSet(global.RedisKeyMessageValues+sendChannel, baseMessageDto.Id, jsonString)
	if nil != err {
		return err
	}

	return global.GetRedisClient().SortedSetAdd(global.RedisKeyMessageKeys+sendChannel, baseMessageDto.Id, float64(baseMessageDto.CreatedTime))
}
