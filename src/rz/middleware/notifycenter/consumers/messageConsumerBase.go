package consumers

import (
	"time"
	"fmt"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/common"
	"runtime"
	"errors"
	"strings"
)

type getMessageFunc func(int, time.Time) (interface{}, *models.PoBase, *models.CallbackBasePo, error)
type sendFunc func(interface{}) (error)
type poToDtoFunc func(interface{}) (interface{})

type messageConsumerBase struct {
	getMessageFunc        getMessageFunc
	sendFunc              sendFunc
	poToDtoFunc           poToDtoFunc
	messageManagementBase *managements.MessageManagementBase
	asyncJobWorker        *common.AsyncJobWorker
	httpClient            *common.HttpClient
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

	return nil
}

func (myself *messageConsumerBase) CloseAndWait() {
	myself.asyncJobWorker.CloseAndWait()
}

func (myself *messageConsumerBase) start(parameter interface{}) (error) {
	now := time.Now()
	messageIds, err := myself.messageManagementBase.DequeueMessageIds(now)
	if nil != err {
		fmt.Printf("failed to get message(%s) ids. error: %s", myself.messageManagementBase.KeySuffix, err.Error())
		return err
	}
	if nil == messageIds {
		return nil
	}

	for _, messageId := range messageIds {
		affectedCount, err := myself.messageManagementBase.RemoveMessageId(messageId)
		if nil != err || 0 == affectedCount {
			fmt.Printf("failed to remove message(%d). error: %s", messageId, err)
			continue
		}

		messagePo, poBase, callbackBasePo, flagError := myself.getMessageFunc(messageId, now)
		if nil == flagError {
			flagError = myself.consume(messagePo, messageId, poBase, callbackBasePo)
		}

		var messageState enumerations.MessageState
		if nil == flagError {
			fmt.Printf("success to consume message(%d)\n", messageId)
			messageState = enumerations.Sent
		} else {
			fmt.Printf("failed to consume message(%d). error: %s\n", messageId, flagError.Error())
			messageState = enumerations.Error
		}
		managements.ModifyMessageFlowAsync(
			myself.messageManagementBase,
			messageId,
			poBase,
			callbackBasePo,
			messageState,
			true,
			time.Now(),
			flagError.Error())

		if "" != callbackBasePo.FinishedCallbackUrls {
			errorMessages := ""
			urls := strings.Split(callbackBasePo.FinishedCallbackUrls, ",")
			for _, url := range urls {
				messageStateCallbackRequestDto := &models.MessageStateCallbackRequestDto{
					Message:      myself.poToDtoFunc(messagePo),
					MessageState: messageState,
				}
				_, err = myself.httpClient.Post(url, messageStateCallbackRequestDto)
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

	err := myself.sendFunc(messagePo)
	if nil != err {
		return err
	}

	return nil
}
