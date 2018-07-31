package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/exceptions"
)

var (
	SmsTemplateService = smsTemplateService{}
)

type smsTemplateService struct {
	MessageServiceBase
}

func (myself *smsTemplateService) Add(smsTemplateDto *models.SmsTemplateDto) (int, error) {
	err := VerifySmsTemplateDto(smsTemplateDto)
	if nil != err {
		return 0, err
	}
	exist, err := managements.SmsTemplateManagement.ExistExtend(smsTemplateDto.Extend)
	if nil != err {
		return 0, err
	}
	if exist {
		return 0, exceptions.ExtendExist().AttachMessage(smsTemplateDto.Extend)
	}

	smsTemplatePo := models.SmsTemplateDtoToPo(smsTemplateDto)

	err = managements.SmsTemplateManagement.Add(smsTemplatePo)
	if nil != err {
		return 0, nil
	}

	return smsTemplatePo.Id, nil
}

func (myself *smsTemplateService) GetAll() ([]*models.SmsTemplateDto, error) {
	keyValues, err := managements.SmsTemplateManagement.GetAll()
	if nil != err {
		return nil, err
	}

	var smsTemplatePos []*models.SmsTemplatePo
	for _, value := range keyValues {
		smsTemplatePos = append(smsTemplatePos, value)
	}

	return models.SmsTemplatePosToDtos(smsTemplatePos), nil
}
