package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
)

var (
	SmsProviderService = smsProviderService{}
)

type smsProviderService struct {
	MessageServiceBase
}

func (myself *smsProviderService) Add(smsProviderDto *models.SmsProviderDto) (string, error) {
	err := VerifySmsProviderDto(smsProviderDto)
	if nil != err {
		return "", err
	}

	smsProviderPo := models.SmsProviderDtoToPo(smsProviderDto)
	err = managements.SmsProviderManagement.Add(smsProviderPo)
	if nil != err {
		return "", err
	}

	return smsProviderPo.Id, nil
}

func (myself *smsProviderService) GetAll() ([]*models.SmsProviderDto, error) {
	keyValues, err := managements.SmsProviderManagement.GetAll()
	if nil != err {
		return nil, err
	}

	var smsProviderPos []*models.SmsProviderPo
	for _, value := range keyValues {
		smsProviderPos = append(smsProviderPos, value)
	}

	return models.SmsProviderPosToDtos(smsProviderPos), nil
}
