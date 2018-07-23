package consumers

import (
	"fmt"
	"errors"
	"rz/middleware/notifycenter/models/external"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/channels"
)

var (
	SmsMessageConsumer *smsMessageConsumer
)

func init() {
	SmsMessageConsumer = &smsMessageConsumer{}
	SmsMessageConsumer.runFunc = SmsMessageConsumer.run
	SmsMessageConsumer.expireRunFunc = SmsMessageConsumer.expireRun
	SmsMessageConsumer.getMessageFunc = SmsMessageConsumer.getMessage
	SmsMessageConsumer.sendFunc = SmsMessageConsumer.send
	SmsMessageConsumer.expireSendFunc = SmsMessageConsumer.expireSend
	SmsMessageConsumer.poToDtoFunc = SmsMessageConsumer.poToDto
	SmsMessageConsumer.messageManagementBase = &managements.SmsMessageManagement.MessageManagementBase
	SmsMessageConsumer.name = SmsMessageConsumer.messageManagementBase.KeySuffix
}

type smsMessageConsumer struct {
	messageConsumerBase
}

func (myself *smsMessageConsumer) send(messagePo interface{}) (error) {
	smsMessagePo, ok := messagePo.(*models.SmsMessagePo)
	err := common.Assert.IsTrueToError(ok, "messagePo.(*models.SmsMessagePo)")
	if nil != err {
		return err
	}

	smsChannel, err := channels.ChooseSmsChannel(smsMessagePo)
	if nil != err {
		return err
	}

	smsMessageResponseExternalDto := &external.SmsMessageResponseExternalDto{}
	err = smsChannel.Do(smsMessagePo, smsMessageResponseExternalDto)
	if nil != err {
		return err
	}

	if 0 != smsMessageResponseExternalDto.ErrorCode {
		message := fmt.Sprintf(
			"ErrorInfo: %s; ActionStatus: %s; ErrorCode: %d",
			smsMessageResponseExternalDto.ErrorInfo,
			smsMessageResponseExternalDto.ActionStatus,
			smsMessageResponseExternalDto.ErrorCode)
		return errors.New(message)
	}
	if 0 != smsMessageResponseExternalDto.Result {
		message := fmt.Sprintf(
			"Errmsg: %s; Result: %d",
			smsMessageResponseExternalDto.Errmsg,
			smsMessageResponseExternalDto.Result)
		return errors.New(message)
	}

	return nil
}

func (myself *smsMessageConsumer) getMessage(messageId int64) (interface{}, *models.PoBase, *models.CallbackBasePo, error) {
	smsMessagePo, err := managements.SmsMessageManagement.GetById(messageId)
	if nil != err {
		return nil, nil, nil, err
	}

	return smsMessagePo, &smsMessagePo.PoBase, &smsMessagePo.CallbackBasePo, nil
}

func (myself *smsMessageConsumer) poToDto(messagePo interface{}) (interface{}) {
	smsMessagePo, ok := messagePo.(*models.SmsMessagePo)
	if !ok {
		return nil
	}

	return models.SmsMessagePoToDto(smsMessagePo)
}

func (myself *smsMessageConsumer) expireSend(messagePo interface{}) (error) {
	smsMessagePo, ok := messagePo.(*models.SmsMessagePo)
	err := common.Assert.IsTrueToError(ok, "messagePo.(*models.SmsMessagePo)")
	if nil != err {
		return err
	}

	smsExpireChannel, err := channels.ChooseSmsExpireChannel(smsMessagePo)
	if nil != err {
		return err
	}

	err = smsExpireChannel.Do(smsMessagePo)
	if nil != err {
		return err
	}

	return nil
}
