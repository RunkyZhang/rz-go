package managements

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/repositories"
	"rz/middleware/notifycenter/enumerations"
	"time"
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

func (myself *smsUserMessageManagement) GetByPhoneNumber(nationCode string, phoneNumber string) ([]models.SmsUserMessagePo, error) {
	return repositories.SmsUserMessageRepository.SelectByPhoneNumber(nationCode, phoneNumber)
}

func (myself *smsUserMessageManagement) GetById(id int, date time.Time) (*models.SmsUserMessagePo, error) {
	return repositories.SmsUserMessageRepository.SelectById(id, date)
}
