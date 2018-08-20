package controllers

import (
	"rz/middleware/notifycenter/services"
	"rz/core/common"
	"rz/middleware/notifycenter/models"
)

// MVC structure
var (
	SmsProviderController = smsProviderController{
		AddSmsProviderControllerPack: &common.ControllerPack{
			Pattern:          "/cloud.appgov.notifycenter.service/smsprovider/add",
			Method:           "POST",
			ControllerFunc:   addSmsProvider,
			ConvertToDtoFunc: ConvertToSmsProviderDto,
		},
		GetAllSmsProvidersControllerPack: &common.ControllerPack{
			Pattern:          "/cloud.appgov.notifycenter.service/smsprovider/getall",
			Method:           "GET",
			ControllerFunc:   getSmsProviders,
			ConvertToDtoFunc: func(body []byte) (interface{}, error) { return nil, nil },
		},
	}
)

type smsProviderController struct {
	ControllerBase

	AddSmsProviderControllerPack  *common.ControllerPack
	GetAllSmsProvidersControllerPack *common.ControllerPack
}

func addSmsProvider(dto interface{}) (interface{}, error) {
	smsProviderDto, ok := dto.(*models.SmsProviderDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.SmsProviderDto)")
	if nil != err {
		return nil, err
	}

	return services.SmsProviderService.Add(smsProviderDto)
}

func getSmsProviders(dto interface{}) (interface{}, error) {
	return services.SmsProviderService.GetAll()
}
