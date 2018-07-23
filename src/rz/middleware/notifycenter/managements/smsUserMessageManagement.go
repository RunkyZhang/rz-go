package managements

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/repositories"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/common"
)

var (
	SmsUserMessageManagement = smsUserMessageManagement{}
)

func init() {
	SmsUserMessageManagement.SendChannel = enumerations.SmsCallback
	SmsUserMessageManagement.KeySuffix, _ = enumerations.SendChannelToString(SmsUserMessageManagement.SendChannel)
	SmsUserMessageManagement.messageRepositoryBase = repositories.SmsUserMessageRepository.MessageRepositoryBase
}

type smsUserMessageManagement struct {
	MessageManagementBase
}

func (myself *smsUserMessageManagement) Add(smsUserMessagePo *models.SmsUserMessagePo) (error) {
	err := common.Assert.IsNotNilToError(smsUserMessagePo, "smsUserMessagePo")
	if nil != err {
		return err
	}
	err = common.Assert.IsNotNilToError(smsUserMessagePo.PoBase, "smsUserMessagePo.PoBase")
	if nil != err {
		return err
	}
	err = common.Assert.IsNotNilToError(smsUserMessagePo.CallbackBasePo, "smsUserMessagePo.CallbackBasePo")
	if nil != err {
		return err
	}

	myself.setPoBase(&smsUserMessagePo.PoBase)
	myself.setCallbackBasePo(&smsUserMessagePo.CallbackBasePo)

	return repositories.SmsUserMessageRepository.Insert(smsUserMessagePo)
}

func (myself *smsUserMessageManagement) ModifySmsMessageId(id int64, smsMessageId int64) (int64, error) {
	return repositories.SmsUserMessageRepository.UpdateSmsMessageIdById(id, smsMessageId)
}

func (myself *smsUserMessageManagement) GetByPhoneNumber(nationCode string, phoneNumber string, year int) ([]*models.SmsUserMessagePo, error) {
	return repositories.SmsUserMessageRepository.SelectByPhoneNumber(nationCode, phoneNumber, year)
}

func (myself *smsUserMessageManagement) GetById(id int64) (*models.SmsUserMessagePo, error) {
	return repositories.SmsUserMessageRepository.SelectById(id)
}

func (myself *smsUserMessageManagement) Query(smsMessageId int64, content string, nationCode string, phoneNumber string, templateId int, year int) ([]*models.SmsUserMessagePo, error) {
	return repositories.SmsUserMessageRepository.Select(smsMessageId, content, nationCode, phoneNumber, templateId, year)
}
