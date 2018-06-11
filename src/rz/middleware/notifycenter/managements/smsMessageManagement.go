package managements

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/repositories"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/common"
)

var (
	SmsMessageManagement = smsMessageManagement{}
)

func init() {
	var err error
	SmsMessageManagement.SendChannel = enumerations.Sms
	SmsMessageManagement.keySuffix, err = enumerations.SendChannelToString(SmsMessageManagement.SendChannel)
	common.Assert.IsNilError(err, "")
}

type smsMessageManagement struct {
	MessageManagementBase

	smsMessageRepository *repositories.SmsMessageRepository
}

func (smsMessageManagement *smsMessageManagement) Add(smsMessagePo *models.SmsMessagePo) (error) {
	return smsMessageManagement.smsMessageRepository.Insert(smsMessagePo)
}

func (smsMessageManagement *smsMessageManagement) GetById(id int) (*models.SmsMessagePo, error) {
	return smsMessageManagement.smsMessageRepository.SelectById(id)
}

func (smsMessageManagement *smsMessageManagement) ModifyById(id int, states string, finished bool, errorMessages string) (int64, error) {
	return smsMessageManagement.smsMessageRepository.UpdateById(id, states, finished, errorMessages)
}
