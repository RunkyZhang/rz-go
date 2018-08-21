package managements

import (
	"rz/middleware/notifycenter/global"
	"rz/core/common"
	"rz/middleware/notifycenter/enumerations"
	"time"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/repositories"
	"fmt"
	"errors"
)

type MessageManagementBase struct {
	managementBase

	SendChannel           enumerations.SendChannel
	KeySuffix             string
	messageRepositoryBase repositories.MessageRepositoryBase
}

func (myself *MessageManagementBase) ModifyStates(id int64, state string, errorMessage string, providerIds string, finished *bool, finishedTime *time.Time, year int) (int64, error) {
	return myself.messageRepositoryBase.UpdateStatesById(id, state, errorMessage, providerIds, finished, finishedTime)
}

func (myself *MessageManagementBase) Disable(id int64) (int64, error) {
	return myself.messageRepositoryBase.UpdateDisableById(id, true)
}

func (myself *MessageManagementBase) RemoveId(id int64) (int64, error) {
	return global.GetRedisClient().SortedSetRemoveByValue(global.RedisKeyMessageIds+myself.KeySuffix, common.Int64ToString(id))
}

func (myself *MessageManagementBase) DequeueIds(now time.Time) ([]int64, error) {
	max := float64(now.Unix())
	messageIds, err := global.GetRedisClient().SortedSetRangeByScore(global.RedisKeyMessageIds+myself.KeySuffix, 0, max)
	if nil != err {
		return nil, err
	}

	var values []int64
	if nil != messageIds {
		for _, messageId := range messageIds {
			value, err := common.StringToInt64(messageId)
			if nil != err {
				continue
			}
			values = append(values, value)
		}
	}

	if 0 == time.Now().Unix()%2 {
		common.SortIntSlice(values)
	} else {
		common.SortReverseIntSlice(values)
	}

	return values, nil
}

func (myself *MessageManagementBase) EnqueueIds(id int64, score int64) (error) {
	return global.GetRedisClient().SortedSetAdd(
		global.RedisKeyMessageIds+myself.KeySuffix,
		common.Int64ToString(id),
		float64(score))
}

func (myself *MessageManagementBase) GenerateId(year int) (int64, error) {
	value, err := global.GetRedisClient().HashIncrement(global.RedisKeyMessageAutoIncrementId, myself.KeySuffix, 1)
	if nil != err {
		return 0, err
	}

	return common.StringToInt64(fmt.Sprintf("%d%d", year, value))
}

func (myself *MessageManagementBase) RemoveExpireId(id int64) (int64, error) {
	return global.GetRedisClient().SortedSetRemoveByValue(global.RedisKeyExpireMessageIds+myself.KeySuffix, common.Int64ToString(id))
}

func (myself *MessageManagementBase) EnqueueExpireIds(id int64, score int64) (error) {
	return global.GetRedisClient().SortedSetAdd(
		global.RedisKeyExpireMessageIds+myself.KeySuffix,
		common.Int64ToString(id),
		float64(score))
}

func (myself *MessageManagementBase) DequeueExpireIds(now time.Time) ([]int64, error) {
	max := float64(now.Unix())
	messageIds, err := global.GetRedisClient().SortedSetRangeByScore(global.RedisKeyExpireMessageIds+myself.KeySuffix, 0, max)
	if nil != err {
		return nil, err
	}

	var values []int64
	if nil != messageIds {
		for _, messageId := range messageIds {
			value, err := common.StringToInt64(messageId)
			if nil != err {
				continue
			}
			values = append(values, value)
		}
	}

	if 0 == time.Now().Unix()%2 {
		common.SortIntSlice(values)
	} else {
		common.SortReverseIntSlice(values)
	}

	return values, nil
}

func (myself *MessageManagementBase) setCallbackBasePo(callbackBasePo *models.CallbackBasePo) {
	callbackBasePo.States = enumerations.MessageStateToString(enumerations.Initial)
	callbackBasePo.Disable = false
	callbackBasePo.Finished = false
}

type MessageFlowJobParameter struct {
	MessageManagementBase *MessageManagementBase
	MessageId             int64
	MessageState          enumerations.MessageState
	ErrorMessage          string
	ProviderIds            string
	Finished              *bool
	FinishedTime          *time.Time
	Year                  int
}

func ModifyMessageFlowAsync(
	messageManagementBase *MessageManagementBase,
	messageId int64,
	currentMessageState enumerations.MessageState,
	messageState enumerations.MessageState,
	errorMessage string,
	providerIds string,
	finished *bool,
	finishedTime *time.Time,
	year int) {
	if "" != errorMessage {
		errorMessage = fmt.Sprintf("(%s)%s", enumerations.MessageStateToString(currentMessageState), errorMessage)
	}
	messageFlowJobParameter := &MessageFlowJobParameter{
		MessageManagementBase: messageManagementBase,
		MessageId:             messageId,
		MessageState:          messageState,
		ErrorMessage:          errorMessage,
		ProviderIds:            providerIds,
		Finished:              finished,
		FinishedTime:          finishedTime,
		Year:                  year,
	}

	asyncJob := &common.AsyncJob{
		Name:
		fmt.Sprintf("%d-%s", messageId, enumerations.MessageStateToString(messageState)),
		Type:      "ModifyMessageFlow",
		RunFunc:   modifyMessageFlow,
		Parameter: messageFlowJobParameter,
	}

	global.AsyncJobWorker.Add(asyncJob)
}

func modifyMessageFlow(parameter interface{}) (error) {
	messageFlowJobParameter, ok := parameter.(*MessageFlowJobParameter)
	err := common.Assert.IsTrueToError(ok, "parameter.(*MessageFlowJobParameter)")
	if nil != err {
		return err
	}
	err = common.Assert.IsTrueToError(nil != messageFlowJobParameter.MessageManagementBase, "nil != messageFlowJobParameter.MessageManagementBase")
	if nil != err {
		return err
	}

	affectedCount, err := messageFlowJobParameter.MessageManagementBase.ModifyStates(
		messageFlowJobParameter.MessageId,
		enumerations.MessageStateToString(messageFlowJobParameter.MessageState),
		messageFlowJobParameter.ErrorMessage,
		messageFlowJobParameter.ProviderIds,
		messageFlowJobParameter.Finished,
		messageFlowJobParameter.FinishedTime,
		messageFlowJobParameter.Year)
	if nil != err || 0 == affectedCount {
		return errors.New(fmt.Sprintf("Failed to modify message(%d) state. error: %s", messageFlowJobParameter.MessageId, err))
	}

	return nil
}
