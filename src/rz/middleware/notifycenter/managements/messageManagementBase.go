package managements

import (
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/enumerations"
	"time"
	"sort"
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

func (myself *MessageManagementBase) ModifyById(id int, states string, finished bool, finishedTime time.Time, errorMessages string, date time.Time) (int64, error) {
	return myself.messageRepositoryBase.UpdateById(id, states, finished, finishedTime, errorMessages, date)
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

func (myself *MessageManagementBase) setCallbackBasePo(callbackBasePo *models.CallbackBasePo) {
	callbackBasePo.States = enumerations.MessageStateToString(enumerations.Initial)
}

type MessageFlowJobParameter struct {
	MessageManagementBase *MessageManagementBase
	MessageId             int
	PoBase                *models.PoBase
	CallbackBasePo        *models.CallbackBasePo
	MessageState          enumerations.MessageState
	Finished              bool
	FinishedTime          time.Time
	ErrorMessage          string
}

func ModifyMessageFlowAsync(
	messageManagementBase *MessageManagementBase,
	messageId int,
	poBase *models.PoBase,
	callbackBasePo *models.CallbackBasePo,
	messageState enumerations.MessageState,
	finished bool,
	finishedTime time.Time,
	errorMessage string) {
	messageFlowJobParameter := &MessageFlowJobParameter{
		MessageManagementBase: messageManagementBase,
		MessageId:             messageId,
		PoBase:                poBase,
		CallbackBasePo:        callbackBasePo,
		MessageState:          messageState,
		Finished:              finished,
		FinishedTime:          finishedTime,
		ErrorMessage:          errorMessage,
	}

	asyncJob := &common.AsyncJob{
		Name:      common.Int32ToString(messageId),
		Type:      "ModifyMessageFlow",
		RunFunc:   modifyMessageFlow,
		Parameter: messageFlowJobParameter,
	}

	global.AsyncWorker.Add(asyncJob)
}

func modifyMessageFlow(parameter interface{}) (error) {
	messageFlowJobParameter, ok := parameter.(*MessageFlowJobParameter)
	err := common.Assert.IsTrueToError(ok, "parameter.(*MessageFlowJobParameter)")
	if nil != err {
		return err
	}
	err = common.Assert.IsNotNilToError(messageFlowJobParameter.MessageManagementBase, "messageFlowJobParameter.messageManagementBase")
	if nil != err {
		return err
	}
	err = common.Assert.IsNotNilToError(messageFlowJobParameter.PoBase, "messageFlowJobParameter.PoBase")
	if nil != err {
		return err
	}
	err = common.Assert.IsNotNilToError(messageFlowJobParameter.CallbackBasePo, "messageFlowJobParameter.CallbackBasePo")
	if nil != err {
		return err
	}

	messageFlowJobParameter.CallbackBasePo.Finished = messageFlowJobParameter.Finished
	messageFlowJobParameter.CallbackBasePo.FinishedTime = messageFlowJobParameter.FinishedTime
	state := enumerations.MessageStateToString(messageFlowJobParameter.MessageState)
	messageFlowJobParameter.CallbackBasePo.States = messageFlowJobParameter.CallbackBasePo.States + "+" + state
	if "" != messageFlowJobParameter.ErrorMessage {
		messageFlowJobParameter.CallbackBasePo.ErrorMessages = messageFlowJobParameter.CallbackBasePo.ErrorMessages + "+++" + messageFlowJobParameter.ErrorMessage
	}
	messageFlowJobParameter.CallbackBasePo.FinishedTime = time.Now()

	affectedCount, err := messageFlowJobParameter.MessageManagementBase.ModifyById(
		messageFlowJobParameter.MessageId,
		messageFlowJobParameter.CallbackBasePo.States,
		messageFlowJobParameter.CallbackBasePo.Finished,
		messageFlowJobParameter.CallbackBasePo.FinishedTime,
		messageFlowJobParameter.CallbackBasePo.ErrorMessages,
		messageFlowJobParameter.PoBase.CreatedTime)
	if nil != err || 0 == affectedCount {
		return errors.New(fmt.Sprintf("failed to modify message(%d) state. error: %s", messageFlowJobParameter.MessageId, err))
	}

	return nil
}
