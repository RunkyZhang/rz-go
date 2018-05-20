package controllers

import (
	"rz/middleware/notifycenter/web"
)

type messageController struct {
	baseController

	SendMailControllerPack *web.ControllerPack
	SendSmsControllerPack  *web.ControllerPack
}

func (messageController *messageController) Enable() {
	messageController.enable(Controller)
}

var (
	Controller = messageController{
		SendMailControllerPack: &web.ControllerPack{
			Pattern:          "/message/send-mail",
			ControllerFunc:   sendMail,
			ConvertToDtoFunc: convertToMailMessageDto,
			VerifyFunc:       verifyMailMessageDto,
		},
		SendSmsControllerPack: &web.ControllerPack{
			Pattern:          "/message/send-sms",
			ControllerFunc:   sendSms,
			ConvertToDtoFunc: convertToSmsMessageDto,
			VerifyFunc:       verifySmsMessageDto,
		},
	}
)

func sendMail(dto interface{}) (interface{}, error) {
	return dto, nil
}

func sendSms(dto interface{}) (interface{}, error) {
	return dto, nil
}
