package controllers

import (
	"rz/middleware/notifycenter/services"
	"rz/middleware/notifycenter/models/external"
	"rz/middleware/notifycenter/common"
)

// MVC structure
var (
	SmsUserCallbackController = smsUserCallbackController{
		SmsUserCallbackControllerPack: &common.ControllerPack{
			Pattern:          "/message/sms-user-callback",
			ControllerFunc:   smsUserCallback,
			ConvertToDtoFunc: ConvertToSmsCallbackMessageDto,
		},
	}
)

type smsUserCallbackController struct {
	ControllerBase

	SmsUserCallbackControllerPack *common.ControllerPack
}

func smsUserCallback(dto interface{}) (interface{}, error) {
	smsUserCallbackMessageRequestExternalDto, ok := dto.(*external.SmsUserCallbackMessageRequestExternalDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*external.SmsUserCallbackMessageRequestExternalDto)")
	if nil != err {
		return nil, err
	}

	return services.SmsUserMessageService.Callback(smsUserCallbackMessageRequestExternalDto), nil
}
