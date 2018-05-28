package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/management"
)

var (
	SmsService = smsService{}
)

type smsService struct {
}

func (*smsService) SendSms(smsMessageDto *models.SmsMessageDto) (string, error) {
	VerifySmsMessageDto(smsMessageDto)

	management.SmsManagement.SendSms(smsMessageDto)

	return "", nil
}
