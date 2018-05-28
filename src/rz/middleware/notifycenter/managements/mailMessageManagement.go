package managements

import (
	"encoding/json"
	"rz/middleware/notifycenter/models"
)

var (
	MailMessageManagement = mailMessageManagement{}
)

type mailMessageManagement struct {
	baseMessageManagement
}

func (mailMessageManagement *mailMessageManagement) AddMailMessage(mailMessageDto *models.MailMessageDto) (error) {
	bytes, err := json.Marshal(mailMessageDto)
	if nil != err {
		return err
	}

	return mailMessageManagement.addMessage(&mailMessageDto.MessageDto, string(bytes))
}
