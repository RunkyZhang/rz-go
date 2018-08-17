package provider

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/core/common"
)

var (
	mailProviders = map[string]*mailProviderBase{}
)

type mailDoFunc func(mailMessagePo *models.MailMessagePo) (error)

func ChooseMailProvider(mailMessagePo *models.MailMessagePo) (*mailProviderBase, error) {
	err := common.Assert.IsTrueToError(nil != mailMessagePo, "nil != mailMessagePo")
	if nil != err {
		return nil, err
	}

	mailProvider, ok := mailProviders[MailDefaultProvider.Id]
	if !ok {
		return nil, exceptions.FailedChooseMailChannel().AttachMessage(mailMessagePo.Id)
	}

	return mailProvider, nil
}

type mailProviderBase struct {
	providerBase

	mailDoFunc mailDoFunc
}

func (myself *mailProviderBase) Do(mailMessagePo *models.MailMessagePo) (error) {
	err := common.Assert.IsTrueToError(nil != mailMessagePo, "nil != mailMessagePo")
	if nil != err {
		return err
	}

	return myself.mailDoFunc(mailMessagePo)
}
