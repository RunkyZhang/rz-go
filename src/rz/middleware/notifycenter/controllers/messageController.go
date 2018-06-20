package controllers

import (
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/services"
	"rz/middleware/notifycenter/models"
)

// MVC structure
var (
	MessageController = messageController{
		SendMailControllerPack: &common.ControllerPack{
			Pattern:          "/message/send-mail",
			ControllerFunc:   sendMail,
			ConvertToDtoFunc: ConvertToMailMessageDto,
		},
		SendSmsControllerPack: &common.ControllerPack{
			Pattern:          "/message/send-sms",
			ControllerFunc:   sendSms,
			ConvertToDtoFunc: ConvertToSmsMessageDto,
		},
	}
)

type messageController struct {
	ControllerBase

	SendMailControllerPack *common.ControllerPack
	SendSmsControllerPack  *common.ControllerPack
}

func sendMail(dto interface{}) (interface{}, error) {
	mailMessageDto, ok := dto.(*models.MailMessageDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.MailMessageDto)")
	if nil != err {
		return nil, err
	}

	return services.MailMessageService.SendMail(mailMessageDto)
}

func sendSms(dto interface{}) (interface{}, error) {
	smsMessageDto, ok := dto.(*models.SmsMessageDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.SmsMessageDto)")
	if nil != err {
		return nil, err
	}

	return services.SmsMessageService.SendSms(smsMessageDto)
}
