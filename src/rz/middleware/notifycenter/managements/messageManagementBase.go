package managements

import (
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/enumerations"
	"time"
	"sort"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/repositories"
)

type MessageManagementBase struct {
	managementBase

	SendChannel           enumerations.SendChannel
	KeySuffix             string
	messageRepositoryBase repositories.MessageRepositoryBase
}

func (myself *MessageManagementBase) ModifyById(id int, states string, finished bool, errorMessages string, date time.Time) (int64, error) {
	return repositories.SmsMessageRepository.UpdateById(id, states, finished, errorMessages, date)
}

func (myself *MessageManagementBase) setCallbackBasePo(callbackBasePo *models.CallbackBasePo) {
	callbackBasePo.States = enumerations.MessageStateToString(enumerations.Initial)
}

func (myself *MessageManagementBase) RemoveMessageId(messageId int) (int64, error) {
	return global.GetRedisClient().SortedSetRemoveByValue(global.RedisKeyMessageIds+myself.KeySuffix, common.Int32ToString(messageId))
}

func (myself *MessageManagementBase) DequeueMessageIds(now time.Time) ([]int, error) {
	max := float64(now.Unix())
	messageIds, err := global.GetRedisClient().SortedSetRangeByScore(global.RedisKeyMessageIds+myself.KeySuffix, 0, max)
	if nil != err {
		return nil, err
	}

	var values []int
	if nil != messageIds {
		for _, messageId := range messageIds {
			value, err := common.StringToInt32(messageId)
			if nil != err {
				continue
			}
			values = append(values, value)
		}
	}

	if 0 == time.Now().Unix()%2 {
		sort.Sort(sort.IntSlice(values))
	} else {
		sort.Sort(sort.Reverse(sort.IntSlice(values)))
	}

	return values, nil
}

func (myself *MessageManagementBase) EnqueueMessageIds(messageId int, score int64) (error) {
	return global.GetRedisClient().SortedSetAdd(
		global.RedisKeyMessageIds+myself.KeySuffix,
		common.Int32ToString(messageId),
		float64(score))
}
