package managements

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/repositories"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/common"
	"time"
)

var (
	SmsMessageManagement = smsMessageManagement{}
)

func init() {
	var err error
	SmsMessageManagement.SendChannel = enumerations.Sms
	SmsMessageManagement.KeySuffix, err = enumerations.SendChannelToString(SmsMessageManagement.SendChannel)
	common.Assert.IsNilError(err, "")
	SmsMessageManagement.messageRepositoryBase = repositories.SmsMessageRepository.MessageRepositoryBase
}

type smsMessageManagement struct {
	MessageManagementBase
}

func (myself *smsMessageManagement) Add(smsMessagePo *models.SmsMessagePo) (error) {
	myself.setPoBase(&smsMessagePo.PoBase)
	myself.setCallbackBasePo(&smsMessagePo.CallbackBasePo)
	smsMessagePo.SendChannel = myself.SendChannel

	return repositories.SmsMessageRepository.Insert(smsMessagePo)
}

func (myself *smsMessageManagement) GetById(id int, date time.Time) (*models.SmsMessagePo, error) {
	return repositories.SmsMessageRepository.SelectById(id, date)
}
