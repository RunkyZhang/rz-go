package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/exceptions"
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
		return false, err
	}

	return true, nil
}

func (myself *systemAliasPermissionService) GetAll() ([]*models.SystemAliasPermissionDto, error) {
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

func (myself *systemAliasPermissionService) Modify(modifySystemAliasPermissionRequestDto *models.ModifySystemAliasPermissionRequestDto) (bool, error) {
	err := VerifyModifySystemAliasPermissionRequestDto(modifySystemAliasPermissionRequestDto)
	if nil != err {
		return false, err
	}

	rowsAffected, err := managements.SystemAliasPermissionManagement.Modify(
		modifySystemAliasPermissionRequestDto.SystemAlias,
		modifySystemAliasPermissionRequestDto.SmsPermission,
		modifySystemAliasPermissionRequestDto.MailPermission,
		modifySystemAliasPermissionRequestDto.SmsDayFrequency,
		modifySystemAliasPermissionRequestDto.SmsHourFrequency,
		modifySystemAliasPermissionRequestDto.SmsMinuteFrequency)

	if 0 == rowsAffected && nil == err {
		return false, exceptions.InvalidSystemAlias().AttachMessage(modifySystemAliasPermissionRequestDto.SystemAlias)
	}

	return 0 < rowsAffected, err
}
