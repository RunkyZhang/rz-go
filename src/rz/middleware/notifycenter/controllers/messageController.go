package controllers

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/web"
	"rz/middleware/notifycenter/services"
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
		SendMailControllerPack: &web.ControllerPack{
			Pattern:          "/message/send-mail",
			ControllerFunc:   sendMail,
			ConvertToDtoFunc: convertToMailMessageDto,
		},
		SendSmsControllerPack: &web.ControllerPack{
			Pattern:          "/message/send-sms",
			ControllerFunc:   sendSms,
			ConvertToDtoFunc: convertToSmsMessageDto,
		},
	}
)

type messageController struct {
	baseController

	SendMailControllerPack *web.ControllerPack
	SendSmsControllerPack  *web.ControllerPack
}

func (messageController *messageController) Enable() {
	messageController.enable(MessageController)
}
