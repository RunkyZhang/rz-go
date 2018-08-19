package provider

import (
	"rz/core/common"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"errors"
)

var (
	smsProviders = map[string]*smsProviderBase{}
)

func ChooseSmsProvider(smsMessagePo *models.SmsMessagePo, excludedIds []string) (*smsProviderBase, error) {
	err := common.Assert.IsTrueToError(nil != smsMessagePo, "nil != smsMessagePo")
	if nil != err {
		return nil, err
	}

	var validSmsProviders []*smsProviderBase
	for _, smsProvider := range smsProviders {
		if nil != excludedIds && 0 != len(excludedIds) {
			for _, excludedId := range excludedIds {
				if smsProvider.Id == excludedId {
					continue
				}
			}
		}

		smsProviderPo, err := managements.SmsProviderManagement.GetById(smsProvider.Id)
		if nil != err {
			continue
		}

		smsProvider.smsProviderPo = smsProviderPo
		if validateSmsProvider(smsMessagePo, smsProviderPo) {
			validSmsProviders = append(validSmsProviders, smsProvider)
		}
	}

	if 0 == len(validSmsProviders){
		return nil,  errors.New("there is no any valid [SmsProvider]")
	}

	//var smsProvider *smsProviderBase
	//var ok bool
	//if 0 == (smsMessagePo.CreatedTime.Second() % 2) {
	//	smsProvider, ok = smsProviders[SmsTencentProvider.Id]
	//} else {
	//	smsProvider, ok = smsProviders[SmsDahanProvider.Id]
	//}
	//

	return smsProvider, nil
}

func validateSmsProvider(smsMessagePo *models.SmsMessagePo, smsProviderPo *models.SmsProviderPo) (bool) {
	if smsProviderPo.Disable {
		return false
	}
	if smsMessagePo.ContentType != (smsMessagePo.ContentType | smsProviderPo.ContentTypes) {
		return false
	}

	return true
}

type smsDoFunc func(smsMessagePo *models.SmsMessagePo, smsTemplatePo *models.SmsTemplatePo) (error)

type smsProviderBase struct {
	providerBase

	smsDoFunc     smsDoFunc
	smsProviderPo *models.SmsProviderPo
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
