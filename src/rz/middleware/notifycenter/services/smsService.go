package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
)

var (
	SmsService = smsService{}
)

type smsService struct {
}

func (*smsService) SendSms(smsMessageDto *models.SmsMessageDto) (string, error) {
	exceptions.VerifySmsMessageDto(smsMessageDto)

	return "", nil
}
