package consumers

import (
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/channels"
)

var (
	SmsUserMessageConsumer *smsUserMessageConsumer
)

func init() {
	SmsUserMessageConsumer = &smsUserMessageConsumer{}
	SmsUserMessageConsumer.runFunc = SmsUserMessageConsumer.run
	SmsUserMessageConsumer.getMessageFunc = SmsUserMessageConsumer.getMessage
	SmsUserMessageConsumer.sendFunc = SmsUserMessageConsumer.send
	SmsUserMessageConsumer.poToDtoFunc = SmsUserMessageConsumer.poToDto
	SmsUserMessageConsumer.messageManagementBase = &managements.SmsUserMessageManagement.MessageManagementBase
	SmsUserMessageConsumer.name = SmsUserMessageConsumer.messageManagementBase.KeySuffix
}

type smsUserMessageConsumer struct {
	messageConsumerBase
}

func (myself *smsUserMessageConsumer) send(messagePo interface{}) (error) {
	smsUserMessagePo, ok := messagePo.(*models.SmsUserMessagePo)
	err := common.Assert.IsTrueToError(ok, "messagePo.(*models.SmsUserMessagePo)")
	if nil != err {
		return err
	}

	smsUserChannel, err := channels.ChooseSmsUserChannel(smsUserMessagePo)
	if nil != err {
		return err
	}

	err = smsUserChannel.Do(smsUserMessagePo)
	if nil != err {
		return err
	}

	return nil
}

func (myself *smsUserMessageConsumer) getMessage(messageId int64) (interface{}, *models.PoBase, *models.CallbackBasePo, error) {
	smsUserMessagePo, err := managements.SmsUserMessageManagement.GetById(messageId)
	if nil != err {
		return nil, nil, nil, err
	}

	return smsUserMessagePo, &smsUserMessagePo.PoBase, &smsUserMessagePo.CallbackBasePo, nil
}

func (myself *smsUserMessageConsumer) poToDto(messagePo interface{}) (interface{}) {
	smsUserMessagePo, ok := messagePo.(*models.SmsUserMessagePo)
	if !ok {
		return nil
	}

	return models.SmsUserMessagePoToDto(smsUserMessagePo)
}
