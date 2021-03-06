package consumers

import (
	"time"
	"strings"

	"rz/middleware/notifycenter/models"
	"rz/core/common"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/global"
	"runtime"
)

type getMessageFunc func(int64) (interface{}, *models.PoBase, *models.CallbackBasePo, error)
type sendFunc func(interface{}) (error)
type poToDtoFunc func(interface{}) (interface{})

type messageConsumerBase struct {
	consumerBase

	getMessageFunc        getMessageFunc
	sendFunc              sendFunc
	poToDtoFunc           poToDtoFunc
	messageManagementBase *managements.MessageManagementBase
	expireAsyncJobTrigger *common.AsyncJobTrigger
	expireRunFunc         common.RunFunc
	expireSendFunc        sendFunc
}

func (myself *messageConsumerBase) Start(duration time.Duration) (error) {
	asyncJob := common.AsyncJob{
		Name:    myself.name,
		Type:    "Consumer Message",
		RunFunc: myself.runFunc,
	}
	myself.asyncJobTrigger = common.NewAsyncJobTrigger(runtime.NumCPU(), duration, asyncJob)
	myself.asyncJobTrigger.Start()

	if nil != myself.expireRunFunc {
		expireAsyncJob := common.AsyncJob{
			Name:    myself.name,
			Type:    "Consumer Expire Message",
			RunFunc: myself.expireRunFunc,
		}
		myself.expireAsyncJobTrigger = common.NewAsyncJobTrigger(1, duration, expireAsyncJob)
		myself.expireAsyncJobTrigger.Start()
	}

	return nil
}

func (myself *messageConsumerBase) run(parameter interface{}) (error) {
	now := time.Now()
	messageIds, err := myself.messageManagementBase.DequeueIds(now)
	if nil != err {
		common.GetLogging().Info(err, "Failed to get message(%s) ids", myself.messageManagementBase.KeySuffix)
		return err
	}
	if nil == messageIds {
		return nil
	}

	for _, messageId := range messageIds {
		affectedCount, err := myself.messageManagementBase.RemoveId(messageId)
		if nil != err {
			common.GetLogging().Warn(err, "Failed to remove message(%d)", messageId)
			continue
		}
		// 0 mean: the other consumer remove it, ignore
		if 0 == affectedCount {
			continue
		}

		messagePo, poBase, callbackBasePo, flagError := myself.getMessageFunc(messageId)
		if nil == flagError {
			flagError = myself.consume(messagePo, messageId, poBase, callbackBasePo)
		}

		var messageState enumerations.MessageState
		var errorMessage string
		if nil == flagError {
			common.GetLogging().Info(nil, "Success to consume message(%d)", messageId)
			messageState = enumerations.Sent
			errorMessage = ""
		} else {
			common.GetLogging().Warn(flagError, "Failed to consume message(%d)", messageId)
			messageState = enumerations.Error
			errorMessage = flagError.Error()
		}
		now := time.Now()
		finished := true
		managements.ModifyMessageFlowAsync(
			myself.messageManagementBase,
			messageId,
			enumerations.Sent,
			messageState,
			errorMessage,
			callbackBasePo.ProviderIds,
			&finished,
			&now,
			poBase.CreatedTime.Year())

		if "" != callbackBasePo.FinishedCallbackUrls {
			errorMessage := ""
			urls := strings.Split(callbackBasePo.FinishedCallbackUrls, ",")
			for _, url := range urls {
				messageStateCallbackRequestDto := &models.MessageStateCallbackRequestDto{
					Message:      myself.poToDtoFunc(messagePo),
					MessageState: messageState,
				}
				_, err = global.HttpClient.Post(url, messageStateCallbackRequestDto)
				if nil != err {
					errorMessage += "[" + exceptions.FailedRequestHttp().AttachError(err).AttachMessage(url).Error() + "]"
				}
			}

			managements.ModifyMessageFlowAsync(
				myself.messageManagementBase,
				messageId,
				enumerations.FinishedSent,
				enumerations.FinishedSent,
				errorMessage,
				"",
				nil,
				nil,
				poBase.CreatedTime.Year())
		}
	}

	return nil
}

func (myself *messageConsumerBase) consume(messagePo interface{}, messageId int64, poBase *models.PoBase, callbackBasePo *models.CallbackBasePo) (error) {
	err := common.Assert.IsTrueToError(nil != messagePo, "nil != messagePo")
	if nil != err {
		return err
	}
	err = common.Assert.IsTrueToError(nil != poBase, "nil != poBase")
	if nil != err {
		return err
	}
	err = common.Assert.IsTrueToError(nil != callbackBasePo, "nil != callbackBasePo")
	if nil != err {
		return err
	}

	managements.ModifyMessageFlowAsync(
		myself.messageManagementBase,
		messageId,
		enumerations.Consuming,
		enumerations.Consuming,
		"",
		"",
		nil,
		nil,
		poBase.CreatedTime.Year())

	if callbackBasePo.Disable {
		return exceptions.MessageDisable().AttachMessage(messageId)
	}

	if time.Now().Unix() > callbackBasePo.ExpireTime.Unix() {
		return exceptions.MessageExpire().AttachMessage(messageId)
	}

	err = myself.sendFunc(messagePo)
	if nil != err {
		return err
	}

	return nil
}

func (myself *messageConsumerBase) expireRun(parameter interface{}) (error) {
	now := time.Now()
	messageIds, err := myself.messageManagementBase.DequeueExpireIds(now)
	if nil != err {
		common.GetLogging().Info(err, "Failed to get expire message(%s) ids", myself.messageManagementBase.KeySuffix)
		return err
	}
	if nil == messageIds {
		return nil
	}

	for _, messageId := range messageIds {
		affectedCount, err := myself.messageManagementBase.RemoveExpireId(messageId)
		if nil != err {
			common.GetLogging().Warn(err, "Failed to remove expire message(%d)", messageId)
			continue
		}
		// 0 mean: the other consumer remove it, ignore
		if 0 == affectedCount {
			continue
		}

		messagePo, poBase, callbackBasePo, flagError := myself.getMessageFunc(messageId)
		if nil == flagError {
			flagError = myself.expireConsume(messagePo, messageId, poBase, callbackBasePo)
		}

		var messageState enumerations.MessageState
		var errorMessage string
		if nil == flagError {
			common.GetLogging().Info(nil, "Success to consume expire message(%d)", messageId)
			messageState = enumerations.ExpireSent
			errorMessage = ""
		} else {
			common.GetLogging().Warn(flagError, "Failed to consume expire message(%d)", messageId)
			messageState = enumerations.ExpireError
			errorMessage = flagError.Error()
		}
		managements.ModifyMessageFlowAsync(
			myself.messageManagementBase,
			messageId,
			enumerations.ExpireSent,
			messageState,
			errorMessage,
			callbackBasePo.ProviderIds,
			nil,
			nil,
			poBase.CreatedTime.Year())
	}

	return nil
}

func (myself *messageConsumerBase) expireConsume(messagePo interface{}, messageId int64, poBase *models.PoBase, callbackBasePo *models.CallbackBasePo) (error) {
	err := common.Assert.IsTrueToError(nil != messagePo, "nil != messagePo")
	if nil != err {
		return err
	}
	err = common.Assert.IsTrueToError(nil != poBase, "nil != poBase")
	if nil != err {
		return err
	}
	err = common.Assert.IsTrueToError(nil != callbackBasePo, "nil != callbackBasePo")
	if nil != err {
		return err
	}

	managements.ModifyMessageFlowAsync(
		myself.messageManagementBase,
		messageId,
		enumerations.ExpireConsuming,
		enumerations.ExpireConsuming,
		"",
		"",
		nil,
		nil,
		poBase.CreatedTime.Year())

	if callbackBasePo.Disable {
		return exceptions.MessageDisable().AttachMessage(messageId)
	}

	err = myself.expireSendFunc(messagePo)
	if nil != err {
		return err
	}

	return nil
}
