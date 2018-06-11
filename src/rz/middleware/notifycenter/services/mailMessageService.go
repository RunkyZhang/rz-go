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
	messageServiceBase
}

func (mailMessageService *mailMessageService) SendMail(mailMessageDto *models.MailMessageDto) (int, error) {
	err := VerifyMailMessageDto(mailMessageDto)
	if nil != err {
		return 0, err
	}

	mailMessagePo := models.MailMessageDtoToPo(mailMessageDto)
	mailMessageService.setMessageBasePo(&mailMessagePo.MessageBasePo)

	err = managements.MailMessageManagement.AddMailMessage(mailMessageDto)
	if nil != err {
		return 0, err
	}

	return mailMessageDto.Id, nil
}
