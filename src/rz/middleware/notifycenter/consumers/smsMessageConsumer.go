package consumers

import (
	"rz/middleware/notifycenter/models"
	"rz/core/common"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/provider"
	"errors"
	"fmt"
	"rz/middleware/notifycenter/enumerations"
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

	var excludedProviderIds []string
	errorMessages := ""
	for ; ; {
		smsProvider, err := provider.ChooseSmsProvider(smsMessagePo, smsTemplatePo, excludedProviderIds)
		if nil != err {
			errorMessage := fmt.Sprintf("The Message(%d) cannot be send by all [SmsProvider]; error: %s", smsMessagePo.Id, err.Error())
			common.GetLogging().Warn(err, errorMessage)
			errorMessages += fmt.Sprintf("[%s]", errorMessage)
			break
		}

		if "" == smsMessagePo.ProviderId {
			smsMessagePo.ProviderId += smsProvider.Id
		} else {
			smsMessagePo.ProviderId += "+" + smsProvider.Id
		}
		err = smsProvider.Do(smsMessagePo, smsTemplatePo)
		if nil == err {
			errorMessages = ""
			break
		} else {
			excludedProviderIds = append(excludedProviderIds, smsProvider.Id)

			errorMessage := fmt.Sprintf("Failed to send the message(%d) with [SmsProvider](%s); error: %s", smsMessagePo.Id, smsProvider.Id, err.Error())
			common.GetLogging().Warn(err, errorMessage)
			errorMessages += fmt.Sprintf("[%s]", errorMessage)
		}

		// Advertisement message only send one time
		if enumerations.SmsContextTypeAdvertisement == smsTemplatePo.ContentType {
			break
		}
	}

	if "" != errorMessages {
		return errors.New(errorMessages[0:(len(errorMessages) - 1)])
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
