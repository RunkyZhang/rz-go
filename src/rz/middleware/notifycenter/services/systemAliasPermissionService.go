package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
)

var (
	SystemAliasPermissionService = systemAliasPermissionService{}
)

type systemAliasPermissionService struct {
	MessageServiceBase
}

func (myself *systemAliasPermissionService) Add(systemAliasPermissionDto *models.SystemAliasPermissionDto) (bool, error) {
	err := VerifySystemAliasPermissionDto(systemAliasPermissionDto)
	if nil != err {
		return false, err
	}

	systemAliasPermissionPo := models.SystemAliasPermissionDtoToPo(systemAliasPermissionDto)
	err = managements.SystemAliasPermissionManagement.Add(systemAliasPermissionPo)
	if nil != err {
		return false, nil
	}

	return true, nil
}

func (myself *systemAliasPermissionService) Get() ([]*models.SystemAliasPermissionDto, error) {
	keyValues, err := managements.SystemAliasPermissionManagement.GetAll()
	if nil != err {
		return nil, err
	}

	var systemAliasPermissionPos []*models.SystemAliasPermissionPo
	for _, value := range keyValues {
		systemAliasPermissionPos = append(systemAliasPermissionPos, value)
	}

	return models.SystemAliasPermissionPosToDtos(systemAliasPermissionPos), nil
}
