package managements

import (
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/common"
)

type messageManagementBase struct {
}

func (*messageManagementBase) addMessage(messageBaseDto *models.MessageBaseDto, jsonString string) (error) {
	sendChannel, err := enumerations.SendChannelToString(messageBaseDto.SendChannel)
	if nil != err {
		return err
	}

	err = global.GetRedisClient().HashSet(global.RedisKeyMessageValues+sendChannel, common.Int32ToString(messageBaseDto.Id), jsonString)
	if nil != err {
		return err
	}

	return global.GetRedisClient().SortedSetAdd(
		global.RedisKeyMessageKeys+sendChannel,
		common.Int32ToString(messageBaseDto.Id),
		float64(messageBaseDto.CreatedTime))
}
