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

type MessageFlowJobParameter struct {
	MessageManagementBase *MessageManagementBase
	MessageId             int
	PoBase                *models.PoBase
	CallbackBasePo        *models.CallbackBasePo
	MessageState          enumerations.MessageState
	Finished              bool
	ErrorMessage          string
}

func ModifyMessageFlowAsync(
	messageManagementBase *MessageManagementBase,
	messageId int,
	poBase *models.PoBase,
	callbackBasePo *models.CallbackBasePo,
	messageState enumerations.MessageState,
	finished bool,
	errorMessage string) {
	messageFlowJobParameter := &MessageFlowJobParameter{
		MessageManagementBase: messageManagementBase,
		MessageId:             messageId,
		PoBase:                poBase,
		CallbackBasePo:        callbackBasePo,
		MessageState:          messageState,
		Finished:              finished,
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
	messageFlowJobParameter := parameter.(*MessageFlowJobParameter)
	err := common.Assert.IsNotNilToError(messageFlowJobParameter, "messageFlowJobParameter")
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
	err = common.Assert.IsNotNilToError(messageFlowJobParameter.CallbackBasePo, "messageFlowJobParameter.messageManagementBase")
	if nil != err {
		return err
	}

	state := enumerations.MessageStateToString(messageFlowJobParameter.MessageState)
	messageFlowJobParameter.CallbackBasePo.States = messageFlowJobParameter.CallbackBasePo.States + "+" + state
	var errorMessages string
	if "" == messageFlowJobParameter.ErrorMessage {
		errorMessages = ""
	} else {
		messageFlowJobParameter.CallbackBasePo.ErrorMessages =
			messageFlowJobParameter.CallbackBasePo.ErrorMessages + "+++" + messageFlowJobParameter.ErrorMessage
		errorMessages = messageFlowJobParameter.CallbackBasePo.ErrorMessages
	}

	affectedCount, err := messageFlowJobParameter.MessageManagementBase.ModifyById(
		messageFlowJobParameter.MessageId,
		messageFlowJobParameter.CallbackBasePo.States,
		messageFlowJobParameter.Finished,
		errorMessages,
		messageFlowJobParameter.PoBase.CreatedTime)
	if nil != err || 0 == affectedCount {
		return errors.New(fmt.Sprintf("failed to modify message(%d) state. error: %s", messageFlowJobParameter.MessageId, err))
	}

	return nil
}
