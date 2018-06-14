package managements

import (
	"rz/middleware/notifycenter/enumerations"
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
	MailMessageManagement.SendChannel = enumerations.Mail
	MailMessageManagement.KeySuffix, _ = enumerations.SendChannelToString(MailMessageManagement.SendChannel)
	MailMessageManagement.messageRepositoryBase = repositories.MailMessageRepository.MessageRepositoryBase
}

func (myself *mailMessageManagement) Add(mailMessagePo *models.MailMessagePo) (error) {
	myself.setPoBase(&mailMessagePo.PoBase)
	myself.setCallbackBasePo(&mailMessagePo.CallbackBasePo)
	mailMessagePo.SendChannel = myself.SendChannel

	return repositories.MailMessageRepository.Insert(mailMessagePo)
}

func (myself *mailMessageManagement) GetById(id int, date time.Time) (*models.MailMessagePo, error) {
	return repositories.MailMessageRepository.SelectById(id, date)
}
