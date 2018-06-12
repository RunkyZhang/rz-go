package managements

import (
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/repositories"
	"rz/middleware/notifycenter/models"
	"time"
)

var (
	MailMessageManagement = mailMessageManagement{}
)

type mailMessageManagement struct {
	MessageManagementBase
}

func init() {
	var err error
	MailMessageManagement.SendChannel = enumerations.Mail
	MailMessageManagement.keySuffix, err = enumerations.SendChannelToString(MailMessageManagement.SendChannel)
	common.Assert.IsNilError(err, "")
}

func (myself *mailMessageManagement) Add(mailMessagePo *models.MailMessagePo) (error) {
	myself.setPoBase(&mailMessagePo.PoBase)

	return repositories.MailMessageRepository.Insert(mailMessagePo)
}

func (myself *mailMessageManagement) GetById(id int, date time.Time) (*models.MailMessagePo, error) {
	return repositories.MailMessageRepository.SelectById(id, date)
}

func (myself *mailMessageManagement) ModifyById(id int, states string, finished bool, errorMessages string, date time.Time) (int64, error) {
	return repositories.MailMessageRepository.UpdateById(id, states, finished, errorMessages, date)
}
