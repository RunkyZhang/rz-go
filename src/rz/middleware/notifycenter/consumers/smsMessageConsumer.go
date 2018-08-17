package consumers

import (
	"rz/middleware/notifycenter/models"
	"rz/core/common"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/provider"
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
	smsTemplatePo, err := managements.SmsTemplateManagement.GetByTemplateId(smsMessagePo.TemplateId)
	if nil != err {
		return err
	}

	smsProvider, err := provider.ChooseSmsProvider(smsMessagePo)
	if nil != err {
		return err
	}

	smsMessagePo.ProviderId = smsProvider.Id
	err = smsProvider.Do(smsMessagePo, smsTemplatePo)
	if nil != err {
		return err
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

	smsExpireProvider, err := provider.ChooseSmsExpireProvider(smsMessagePo)
	if nil != err {
		return err
	}

	smsMessagePo.ProviderId += "+" + smsExpireProvider.Id
	err = smsExpireProvider.Do(smsMessagePo)
	if nil != err {
		return err
	}

	return nil
}
