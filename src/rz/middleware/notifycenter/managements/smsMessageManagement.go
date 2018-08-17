package managements

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/repositories"
	"rz/middleware/notifycenter/enumerations"
	"rz/core/common"
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
	err := common.Assert.IsTrueToError(nil != smsMessagePo, "nil != smsMessagePo")
	if nil != err {
		return err
	}

	myself.setPoBase(&smsMessagePo.PoBase)
	myself.setCallbackBasePo(&smsMessagePo.CallbackBasePo)
	smsMessagePo.SendChannel = myself.SendChannel

	return repositories.SmsMessageRepository.Insert(smsMessagePo)
}

func (myself *smsMessageManagement) GetById(id int64) (*models.SmsMessagePo, error) {
	return repositories.SmsMessageRepository.SelectById(id)
}

func (myself *smsMessageManagement) GetByIds(id []int64, year int) ([]*models.SmsMessagePo, error) {
	return repositories.SmsMessageRepository.SelectByIds(id, year)
}

func (myself *smsMessageManagement) GetByIdentifyingCode(templateId int, identifyingCode string, year int) (*models.SmsMessagePo, error) {
	return repositories.SmsMessageRepository.SelectByIdentifyingCode(templateId, identifyingCode, year)
}
