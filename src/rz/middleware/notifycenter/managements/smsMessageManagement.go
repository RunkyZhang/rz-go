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
	SmsMessageManagement.keySuffix, err = enumerations.SendChannelToString(SmsMessageManagement.SendChannel)
	common.Assert.IsNilError(err, "")
}

type smsMessageManagement struct {
	MessageManagementBase
}

func (myself *smsMessageManagement) Add(smsMessagePo *models.SmsMessagePo) (error) {
	myself.setPoBase(&smsMessagePo.PoBase)

	return repositories.SmsMessageRepository.Insert(smsMessagePo)
}

func (myself *smsMessageManagement) GetById(id int, date time.Time) (*models.SmsMessagePo, error) {
	return repositories.SmsMessageRepository.SelectById(id, date)
}

func (myself *smsMessageManagement) ModifyById(id int, states string, finished bool, errorMessages string, date time.Time) (int64, error) {
	return repositories.SmsMessageRepository.UpdateById(id, states, finished, errorMessages, date)
}
