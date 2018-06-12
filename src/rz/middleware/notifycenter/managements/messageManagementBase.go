package managements

import (
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/enumerations"
	"time"
	"sort"
)

type MessageManagementBase struct {
	managementBase

	SendChannel enumerations.SendChannel
	keySuffix   string
}

func (myself *MessageManagementBase) RemoveMessageId(messageId int) (int64, error) {
	return global.GetRedisClient().SortedSetRemoveByValue(global.RedisKeyMessageIds+myself.keySuffix, common.Int32ToString(messageId))
}

func (myself *MessageManagementBase) DequeueMessageIds(now time.Time) ([]int, error) {
	max := float64(now.Unix())
	messageIds, err := global.GetRedisClient().SortedSetRangeByScore(global.RedisKeyMessageIds+myself.keySuffix, 0, max)
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
