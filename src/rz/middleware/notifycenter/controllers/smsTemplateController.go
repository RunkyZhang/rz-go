package controllers

import (
	"rz/middleware/notifycenter/services"
	"rz/core/common"
	"rz/middleware/notifycenter/models"
)

// MVC structure
var (
	SmsTemplateController = smsTemplateController{
		AddSmsTemplateControllerPack: &common.ControllerPack{
			Pattern:          "/cloud.appgov.notifycenter.service/smstemplate/add",
			Method:           "POST",
			ControllerFunc:   addSmsTemplate,
			ConvertToDtoFunc: ConvertToSmsTemplateDto,
		},
		GetAllSmsTemplatesControllerPack: &common.ControllerPack{
			Pattern:          "/cloud.appgov.notifycenter.service/smstemplate/getall",
			Method:           "GET",
			ControllerFunc:   getSmsTemplates,
			ConvertToDtoFunc: func(body []byte) (interface{}, error) { return nil, nil },
		},
	}
)

type smsTemplateController struct {
	ControllerBase

	AddSmsTemplateControllerPack  *common.ControllerPack
	GetAllSmsTemplatesControllerPack *common.ControllerPack
}

func addSmsTemplate(dto interface{}) (interface{}, error) {
	smsTemplateDto, ok := dto.(*models.SmsTemplateDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.SmsTemplateDto)")
	if nil != err {
		return nil, err
	}

	return services.SmsTemplateService.Add(smsTemplateDto)
}

func getSmsTemplates(dto interface{}) (interface{}, error) {
	return services.SmsTemplateService.GetAll()
}
