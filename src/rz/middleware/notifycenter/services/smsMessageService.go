package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/enumerations"
)

var (
	SmsMessageService = smsMessageService{}
)

func init() {
	MailMessageService.SendChannel = enumerations.Sms
	MailMessageService.Prefix = "S"
}

type smsMessageService struct {
	baseMessageService
}

func (smsMessageService *smsMessageService) SendSms(smsMessageDto *models.SmsMessageDto) (string, error) {
	err := VerifySmsMessageDto(smsMessageDto)
	if nil != err {
		return "", err
	}

	err = smsMessageService.setMessageDto(&smsMessageDto.MessageDto)
	if nil != err {
		return "", err
	}

	err = managements.SmsMessageManagement.AddSmsMessage(smsMessageDto)
	if nil != err {
		return "", err
	}

	return smsMessageDto.Id, nil
}
