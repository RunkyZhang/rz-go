package controllers

import (
	"rz/middleware/notifycenter/services"
	"rz/middleware/notifycenter/models/external"
	"rz/middleware/notifycenter/common"
)

func smsUserCallback(dto interface{}) (interface{}, error) {
	smsUserCallbackMessageRequestExternalDto := dto.(*external.SmsUserCallbackMessageRequestExternalDto)

	return services.SmsUserMessageService.Add(smsUserCallbackMessageRequestExternalDto), nil
}

// MVC structure
var (
	SmsUserCallbackController = smsUserCallbackController{
		SmsUserCallbackControllerPack: &common.ControllerPack{
			Pattern:          "/message/sms-callback",
			ControllerFunc:   smsUserCallback,
			ConvertToDtoFunc: ConvertToSmsCallbackMessageDto,
		},
	}
)

type smsUserCallbackController struct {
	ControllerBase

	SmsUserCallbackControllerPack *common.ControllerPack
}

func (smsCallbackController *smsUserCallbackController) Enable() {
	smsCallbackController.ControllerBase.Enable(MessageController)
}
