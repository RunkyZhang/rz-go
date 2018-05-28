package services

import (
	"rz/middleware/notifycenter/models"
)

var (
	MailService = mailService{}
)

type mailService struct {
}

func (*mailService) SendMail(mailMessageDto *models.MailMessageDto) (string, error) {
	VerifyMailMessageDto(mailMessageDto)



	return mailMessageDto.Tos[0], nil
}
