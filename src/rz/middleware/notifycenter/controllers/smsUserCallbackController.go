package controllers

import (
	"rz/middleware/notifycenter/web"
	"rz/middleware/notifycenter/services"
	"rz/middleware/notifycenter/models/external"
)

func smsUserCallback(dto interface{}) (interface{}, error) {
	smsUserCallbackMessageRequestExternalDto := dto.(*external.SmsUserCallbackMessageRequestExternalDto)

	return services.SmsUserCallbackService.Add(smsUserCallbackMessageRequestExternalDto), nil
}

// MVC structure
var (
	SmsUserCallbackController = smsUserCallbackController{
		SmsUserCallbackControllerPack: &web.ControllerPack{
			Pattern:          "/message/sms-callback",
			ControllerFunc:   smsUserCallback,
			ConvertToDtoFunc: ConvertToSmsCallbackMessageDto,
		},
	}
)

type smsUserCallbackController struct {
	ControllerBase

	SmsUserCallbackControllerPack *web.ControllerPack
}

func (smsCallbackController *smsUserCallbackController) Enable() {
	smsCallbackController.ControllerBase.Enable(MessageController)
}
