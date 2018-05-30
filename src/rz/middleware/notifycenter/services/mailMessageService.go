package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/managements"
)

var (
	MailMessageService = mailMessageService{}
)

func init() {
	MailMessageService.SendChannel = enumerations.Mail
	MailMessageService.Prefix = "E"
}

type mailMessageService struct {
	baseMessageService
}

func (mailMessageService *mailMessageService) SendMail(mailMessageDto *models.MailMessageDto) (string, error) {
	err := VerifyMailMessageDto(mailMessageDto)
	if nil != err {
		return "", err
	}

	err = mailMessageService.setMessageDto(&mailMessageDto.BaseMessageDto)
	if nil != err {
		return "", err
	}

	err = managements.MailMessageManagement.AddMailMessage(mailMessageDto)
	if nil != err {
		return "", err
	}

	return mailMessageDto.Id, nil
}
