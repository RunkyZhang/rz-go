package provider

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/core/common"
)

var (
	smsExpireProviders = map[string]*smsExpireProviderBase{}
)

type smsExpireDoFunc func(smsMessagePo *models.SmsMessagePo) (error)

func ChooseSmsExpireProvider(smsMessagePo *models.SmsMessagePo) (*smsExpireProviderBase, error) {
	err := common.Assert.IsTrueToError(nil != smsMessagePo, "nil != smsMessagePo")
	if nil != err {
		return nil, err
	}

	smsExpireProvider, ok := smsExpireProviders[SmsExpireDefaultProvider.Id]
	if !ok {
		return nil, exceptions.FailedChooseSmsExpireChannel().AttachMessage(smsMessagePo.Id)
	}

	return smsExpireProvider, nil
}

type smsExpireProviderBase struct {
	providerBase

	smsExpireDoFunc smsExpireDoFunc
}

func (myself *smsExpireProviderBase) Do(smsMessagePo *models.SmsMessagePo) (error) {
	err := common.Assert.IsTrueToError(nil != smsMessagePo, "nil != smsMessagePo")
	if nil != err {
		return err
	}

	return myself.smsExpireDoFunc(smsMessagePo)
}
