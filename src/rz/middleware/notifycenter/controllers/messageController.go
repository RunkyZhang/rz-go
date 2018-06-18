package controllers

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/services"
	"rz/middleware/notifycenter/common"
)

func sendMail(dto interface{}) (interface{}, error) {
	mailMessageDto := dto.(*models.MailMessageDto)

	return services.MailMessageService.SendMail(mailMessageDto)
}

func sendSms(dto interface{}) (interface{}, error) {
	smsMessageDto := dto.(*models.SmsMessageDto)

	return services.SmsMessageService.SendSms(smsMessageDto)
}

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

func (messageController *messageController) Enable() {
	messageController.ControllerBase.Enable(MessageController)
}
