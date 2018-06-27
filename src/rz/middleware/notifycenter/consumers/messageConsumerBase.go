package consumers

import (
	"time"
	"fmt"
	"runtime"
	"errors"
	"strings"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/exceptions"
)

type getMessageFunc func(int, time.Time) (interface{}, *models.PoBase, *models.CallbackBasePo, error)
type sendFunc func(interface{}) (error)
type poToDtoFunc func(interface{}) (interface{})

var httpClient = common.NewHttpClient(nil)

type messageConsumerBase struct {
	getMessageFunc        getMessageFunc
	sendFunc              sendFunc
	poToDtoFunc           poToDtoFunc
	messageManagementBase *managements.MessageManagementBase
	asyncJobWorker        *common.AsyncJobWorker
}

func (myself *messageConsumerBase) Start(duration time.Duration) (error) {
	if nil == myself.messageManagementBase {
		return errors.New("[myself.messageManagementBase] is nil")
	}

	asyncJob := &common.AsyncJob{
		Name:    myself.messageManagementBase.KeySuffix,
		Type:    "Consumer",
		RunFunc: myself.start,
	}
	myself.asyncJobWorker = common.NewAsyncJobWorker(runtime.NumCPU(), duration, asyncJob)
	myself.asyncJobWorker.Start()

	return nil
}

func (myself *messageConsumerBase) CloseAndWait() {
	myself.asyncJobWorker.CloseAndWait()
}

func (myself *messageConsumerBase) start(parameter interface{}) (error) {
	now := time.Now()
	messageIds, err := myself.messageManagementBase.DequeueMessageIds(now)
	if nil != err {
		common.GetLogging().Info(err, "failed to get message(%s) ids", myself.messageManagementBase.KeySuffix)
		return err
	}
	if nil == messageIds {
		return nil
	}

	for _, messageId := range messageIds {
		affectedCount, err := myself.messageManagementBase.RemoveMessageId(messageId)
		if nil != err {
			common.GetLogging().Error(err, "failed to remove message(%d)", messageId)
			continue
		}
		// 0 mean: the other consumer remove it, ignore
		if 0 == affectedCount {
			continue
		}

		messagePo, poBase, callbackBasePo, flagError := myself.getMessageFunc(messageId, now)
		if nil == flagError {
			flagError = myself.consume(messagePo, messageId, poBase, callbackBasePo)
		}

		var messageState enumerations.MessageState
		var errorMessage string
		if nil == flagError {
			common.GetLogging().Info(nil, "success to consume message(%d)", messageId)
			messageState = enumerations.Sent
			errorMessage = ""
		} else {
			common.GetLogging().Error(flagError, "failed to consume message(%d)", messageId)
			messageState = enumerations.Error
			errorMessage = flagError.Error()
		}
		managements.ModifyMessageFlowAsync(
			myself.messageManagementBase,
			messageId,
			poBase,
			callbackBasePo,
			messageState,
			true,
			time.Now(),
			errorMessage)

		if "" != callbackBasePo.FinishedCallbackUrls {
			errorMessages := ""
			urls := strings.Split(callbackBasePo.FinishedCallbackUrls, ",")
			for _, url := range urls {
				messageStateCallbackRequestDto := &models.MessageStateCallbackRequestDto{
					Message:      myself.poToDtoFunc(messagePo),
					MessageState: messageState,
				}
				_, err = httpClient.Post(url, messageStateCallbackRequestDto)
				if nil != err {
					errorMessages += errorMessages + fmt.Sprintf("+++failed to invoke url(%s)", url)
				}
			}

			managements.ModifyMessageFlowAsync(
				myself.messageManagementBase,
				messageId,
				poBase,
				callbackBasePo,
				enumerations.FinishedCallbackInvoked,
				true,
				callbackBasePo.FinishedTime,
				errorMessages)
		}
	}

	return nil
}

func (myself *messageConsumerBase) consume(messagePo interface{}, messageId int, poBase *models.PoBase, callbackBasePo *models.CallbackBasePo) (error) {
	err := common.Assert.IsNotNilToError(messagePo, "messagePo")
	if nil != err {
		return err
	}
	err = common.Assert.IsNotNilToError(poBase, "poBase")
	if nil != err {
		return err
	}
	err = common.Assert.IsNotNilToError(callbackBasePo, "callbackBasePo")
	if nil != err {
		return err
	}

	managements.ModifyMessageFlowAsync(
		myself.messageManagementBase,
		messageId,
		poBase,
		callbackBasePo,
		enumerations.Consuming,
		false,
		time.Unix(0, 0),
		"")

	if time.Now().Unix() > callbackBasePo.ExpireTime.Unix() {
		return exceptions.MessageExpire().AttachMessage(common.Int32ToString(messageId))
	}

	err = myself.sendFunc(messagePo)
	if nil != err {
		return err
	}

	return nil
}
