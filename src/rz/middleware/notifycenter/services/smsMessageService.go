package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/enumerations"
	"fmt"
)

var (
	SmsMessageService = smsMessageService{}
)

func init() {
	SmsMessageService.SendChannel = enumerations.Sms
	SmsMessageService.Prefix = "S"
}

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
	fmt.Println(smsMessagePo.Id)
	if nil != err {
		return 0, err
	}

	return smsMessageDto.Id, nil
}
