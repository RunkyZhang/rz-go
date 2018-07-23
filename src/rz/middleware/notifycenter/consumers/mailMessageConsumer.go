package consumers

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/channels"
)

var (
	MailMessageConsumer *mailMessageConsumer
)

func init() {
	MailMessageConsumer = &mailMessageConsumer{}
	MailMessageConsumer.runFunc = MailMessageConsumer.run
	MailMessageConsumer.getMessageFunc = MailMessageConsumer.getMessage
	MailMessageConsumer.sendFunc = MailMessageConsumer.send
	MailMessageConsumer.poToDtoFunc = MailMessageConsumer.poToDto
	MailMessageConsumer.messageManagementBase = &managements.MailMessageManagement.MessageManagementBase
	MailMessageConsumer.name = MailMessageConsumer.messageManagementBase.KeySuffix
}

type mailMessageConsumer struct {
	messageConsumerBase
}

func (myself *mailMessageConsumer) send(messagePo interface{}) error {
	mailMessagePo, ok := messagePo.(*models.MailMessagePo)
	err := common.Assert.IsTrueToError(ok, "messagePo.(*models.MailMessagePo)")
	if nil != err {
		return err
	}

	mailChannel, err := channels.ChooseMailChannel(mailMessagePo)
	if nil != err {
		return err
	}

	err = mailChannel.Do(mailMessagePo)
	if nil != err {
		return err
	}

	return nil
}

func (myself *mailMessageConsumer) getMessage(messageId int64) (interface{}, *models.PoBase, *models.CallbackBasePo, error) {
	mailMessagePo, err := managements.MailMessageManagement.GetById(messageId)
	if nil != err {
		return nil, nil, nil, err
	}

	return mailMessagePo, &mailMessagePo.PoBase, &mailMessagePo.CallbackBasePo, nil
}

func (myself *mailMessageConsumer) poToDto(messagePo interface{}) (interface{}) {
	mailMessagePo, ok := messagePo.(*models.MailMessagePo)
	if !ok {
		return nil
	}

	return models.MailMessagePoToDto(mailMessagePo)
}
