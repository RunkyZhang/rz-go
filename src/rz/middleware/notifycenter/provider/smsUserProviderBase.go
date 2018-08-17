package provider

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/core/common"
)

var (
	smsUserProviders = map[string]*smsUserProviderBase{}
)

type smsUserDoFunc func(smsUserMessagePo *models.SmsUserMessagePo) (error)

func ChooseSmsUserProvider(smsUserMessagePo *models.SmsUserMessagePo) (*smsUserProviderBase, error) {
	err := common.Assert.IsTrueToError(nil != smsUserMessagePo, "nil != smsUserMessagePo")
	if nil != err {
		return nil, err
	}

	smsUserProvider, ok := smsUserProviders[SmsUserDefaultProvider.Id]
	if !ok {
		return nil, exceptions.FailedChooseSmsUserChannel().AttachMessage(smsUserMessagePo.Id)
	}

	return smsUserProvider, nil
}

type smsUserProviderBase struct {
	providerBase

	smsUserDoFunc smsUserDoFunc
}

func (myself *smsUserProviderBase) Do(smsUserMessagePo *models.SmsUserMessagePo) (error) {
	err := common.Assert.IsTrueToError(nil != smsUserMessagePo, "nil != smsUserMessagePo")
	if nil != err {
		return err
	}

	return myself.smsUserDoFunc(smsUserMessagePo)
}
