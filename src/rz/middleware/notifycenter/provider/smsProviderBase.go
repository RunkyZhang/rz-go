package provider

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/core/common"
)

var (
	smsProviders = map[string]*smsProviderBase{}
)

type smsDoFunc func(smsMessagePo *models.SmsMessagePo, smsTemplatePo *models.SmsTemplatePo) (error)

func ChooseSmsProvider(smsMessagePo *models.SmsMessagePo) (*smsProviderBase, error) {
	err := common.Assert.IsTrueToError(nil != smsMessagePo, "nil != smsMessagePo")
	if nil != err {
		return nil, err
	}

	var smsProvider *smsProviderBase
	var ok bool
	if 0 == (smsMessagePo.CreatedTime.Second() % 2) {
		smsProvider, ok = smsProviders[SmsTencentProvider.Id]
	} else {
		smsProvider, ok = smsProviders[SmsDahanProvider.Id]
	}

	if !ok {
		return nil, exceptions.FailedChooseSmsChannel().AttachMessage(smsMessagePo.Id)
	}

	return smsProvider, nil
}

type smsProviderBase struct {
	providerBase

	smsDoFunc smsDoFunc
}

func (myself *smsProviderBase) Do(smsMessagePo *models.SmsMessagePo, smsTemplatePo *models.SmsTemplatePo) (error) {
	err := common.Assert.IsTrueToError(nil != smsMessagePo, "nil != smsMessagePo")
	if nil != err {
		return err
	}
	err = common.Assert.IsTrueToError(nil != smsTemplatePo, "nil != smsTemplatePo")
	if nil != err {
		return err
	}

	return myself.smsDoFunc(smsMessagePo, smsTemplatePo)
}
