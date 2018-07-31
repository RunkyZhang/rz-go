package controllers

import (
	"rz/middleware/notifycenter/services"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/models"
)

// MVC structure
var (
	SystemAliasPermissionController = systemAliasPermissionController{
		AddSystemAliasPermissionControllerPack: &common.ControllerPack{
			Pattern:          "/cloud.appgov.notifycenter.service/permission/add",
			Method:           "POST",
			ControllerFunc:   addSystemAliasPermission,
			ConvertToDtoFunc: ConvertToSystemAliasPermissionDto,
		},
		ModifySystemAliasPermissionControllerPack: &common.ControllerPack{
			Pattern:          "/cloud.appgov.notifycenter.service/permission/modify",
			Method:           "POST",
			ControllerFunc:   modifySystemAliasPermission,
			ConvertToDtoFunc: ConvertToModifySystemAliasPermissionRequestDto,
		},
		GetAllSystemAliasPermissionsControllerPack: &common.ControllerPack{
			Pattern:          "/cloud.appgov.notifycenter.service/permission/getall",
			Method:           "GET",
			ControllerFunc:   getSystemAliasPermissions,
			ConvertToDtoFunc: func(body []byte) (interface{}, error) { return nil, nil },
		},
	}
)

type systemAliasPermissionController struct {
	ControllerBase

	AddSystemAliasPermissionControllerPack     *common.ControllerPack
	ModifySystemAliasPermissionControllerPack  *common.ControllerPack
	GetAllSystemAliasPermissionsControllerPack *common.ControllerPack
}

func addSystemAliasPermission(dto interface{}) (interface{}, error) {
	systemAliasPermissionDto, ok := dto.(*models.SystemAliasPermissionDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.SystemAliasPermissionDto)")
	if nil != err {
		return nil, err
	}

	return services.SystemAliasPermissionService.Add(systemAliasPermissionDto)
}

func modifySystemAliasPermission(dto interface{}) (interface{}, error) {
	modifySystemAliasPermissionRequestDto, ok := dto.(*models.ModifySystemAliasPermissionRequestDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.ModifySystemAliasPermissionRequestDto)")
	if nil != err {
		return nil, err
	}

	return services.SystemAliasPermissionService.Modify(modifySystemAliasPermissionRequestDto)
}

func getSystemAliasPermissions(dto interface{}) (interface{}, error) {
	return services.SystemAliasPermissionService.GetAll()
}
