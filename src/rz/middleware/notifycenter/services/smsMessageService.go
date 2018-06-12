package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
)

var (
	SmsMessageService = smsMessageService{}
)
type smsMessageService struct {
	messageServiceBase
}

func (smsMessageService *smsMessageService) SendSms(smsMessageDto *models.SmsMessageDto) (int, error) {
	err := VerifySmsMessageDto(smsMessageDto)
	if nil != err {
		return 0, err
	}

	smsMessagePo := models.SmsMessageDtoToPo(smsMessageDto)
	smsMessageService.setMessageBasePo(&smsMessagePo.MessageBasePo)

	err = managements.SmsMessageManagement.Add(smsMessagePo)
	if nil != err {
		return 0, err
	}

	return smsMessagePo.Id, nil
}
