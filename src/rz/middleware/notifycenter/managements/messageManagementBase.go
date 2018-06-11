package managements

import (
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/enumerations"
	"time"
	"sort"
)

type MessageManagementBase struct {
	SendChannel enumerations.SendChannel
	keySuffix   string
}

func (messageManagementBase *MessageManagementBase) RemoveMessageId(messageId int) (int64, error) {
	return global.GetRedisClient().SortedSetRemoveByValue(global.RedisKeyMessageIds+messageManagementBase.keySuffix, common.Int32ToString(messageId))
}

func (messageManagementBase *MessageManagementBase) DequeueMessageIds() ([]int, error) {
	max := float64(time.Now().Unix())
	messageIds, err := global.GetRedisClient().SortedSetRangeByScore(global.RedisKeyMessageIds+messageManagementBase.keySuffix, 0, max)
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
