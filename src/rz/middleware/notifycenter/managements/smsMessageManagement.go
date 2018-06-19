package managements

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/repositories"
	"rz/middleware/notifycenter/enumerations"
	"time"
	"rz/middleware/notifycenter/common"
)

var (
	SmsMessageManagement = smsMessageManagement{}
)

func init() {
	SmsMessageManagement.SendChannel = enumerations.Sms
	SmsMessageManagement.KeySuffix, _ = enumerations.SendChannelToString(SmsMessageManagement.SendChannel)
	SmsMessageManagement.messageRepositoryBase = repositories.SmsMessageRepository.MessageRepositoryBase
}

type smsMessageManagement struct {
	MessageManagementBase
}

func (myself *smsMessageManagement) Add(smsMessagePo *models.SmsMessagePo) (error) {
	err := common.Assert.IsNotNilToError(smsMessagePo, "smsMessagePo")
	if nil != err {
		return err
	}
	err = common.Assert.IsNotNilToError(smsMessagePo.PoBase, "smsMessagePo.PoBase")
	if nil != err {
		return err
	}
	err = common.Assert.IsNotNilToError(smsMessagePo.CallbackBasePo, "smsMessagePo.CallbackBasePo")
	if nil != err {
		return err
	}

	myself.setPoBase(&smsMessagePo.PoBase)
	myself.setCallbackBasePo(&smsMessagePo.CallbackBasePo)
	smsMessagePo.SendChannel = myself.SendChannel

	return repositories.SmsMessageRepository.Insert(smsMessagePo)
}

func (myself *smsMessageManagement) GetById(id int, date time.Time) (*models.SmsMessagePo, error) {
	return repositories.SmsMessageRepository.SelectById(id, date)
}

func (myself *smsMessageManagement) GetByIdentifyingCode(templateId int, identifyingCode string, date time.Time) (*models.SmsMessagePo, error) {
	return repositories.SmsMessageRepository.SelectByIdentifyingCode(templateId, identifyingCode, date)
}
