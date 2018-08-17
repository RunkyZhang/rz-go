package controllers

import (
	"rz/middleware/notifycenter/services"
	"rz/middleware/notifycenter/models/external"
	"rz/core/common"
)

// MVC structure
var (
	SmsUserCallbackController = smsUserCallbackController{
		SmsUserCallbackControllerPack: &common.ControllerPack{
			Pattern:          "/message/sms-user-callback",
			Method:           "POST",
			ControllerFunc:   tencentSmsUserCallback,
			ConvertToDtoFunc: ConvertToTencentSmsUserCallbackRequestDto,
		},
		TencentSmsUserCallbackControllerPack: &common.ControllerPack{
			Pattern:          "/message/tencent-user-sms-callback",
			Method:           "POST",
			ControllerFunc:   tencentSmsUserCallback,
			ConvertToDtoFunc: ConvertToTencentSmsUserCallbackRequestDto,
		},
		DahanSmsUserCallbackControllerPack: &common.ControllerPack{
			Pattern:          "/message/dahan-user-sms-callback",
			Method:           "POST",
			ControllerFunc:   dahanSmsUserCallback,
			ConvertToDtoFunc: ConvertToDahanSmsUserCallbackRequestDto,
		},
	}
)

type smsUserCallbackController struct {
	ControllerBase

	SmsUserCallbackControllerPack        *common.ControllerPack
	TencentSmsUserCallbackControllerPack *common.ControllerPack
	DahanSmsUserCallbackControllerPack   *common.ControllerPack
}

func tencentSmsUserCallback(dto interface{}) (interface{}, error) {
	tencentSmsUserCallbackRequestDto, ok := dto.(*external.TencentSmsUserCallbackRequestDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*external.TencentSmsUserCallbackRequestDto)")
	if nil != err {
		return nil, err
	}

	return services.SmsUserMessageService.TencentCallback(tencentSmsUserCallbackRequestDto), nil
}

func dahanSmsUserCallback(dto interface{}) (interface{}, error) {
	dahanSmsUserCallbackRequestDto, ok := dto.(*external.DahanSmsUserCallbackRequestDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*external.DahanSmsUserCallbackRequestDto)")
	if nil != err {
		return nil, err
	}

	return services.SmsUserMessageService.DahanCallbacks(dahanSmsUserCallbackRequestDto), nil
}
