package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
)

var (
	MailService = mailService{}
)

type mailService struct {
}

func (*mailService) SendMail(mailMessageDto *models.MailMessageDto) (string, error) {
	exceptions.VerifyMailMessageDto(mailMessageDto)



	return mailMessageDto.Tos[0], nil
}
