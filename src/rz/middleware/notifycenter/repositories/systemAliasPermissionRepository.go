package repositories

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/common"
)

var (
	SystemAliasPermissionRepository systemAliasPermissionRepository
)

func init() {
	SystemAliasPermissionRepository.DefaultDatabaseKey = "default"
	SystemAliasPermissionRepository.RawTableName = "systemAliasPermissionPo"
}

type systemAliasPermissionRepository struct {
	common.RepositoryBase
}

func (myself *systemAliasPermissionRepository) Insert(systemAliasPermissionPo *models.SystemAliasPermissionPo) (error) {
	err := common.Assert.IsNotNilToError(systemAliasPermissionPo, "systemAliasPermissionPo")
	if nil != err {
		return err
	}

	return myself.RepositoryBase.Insert(systemAliasPermissionPo, nil)
}

func (myself *systemAliasPermissionRepository) SelectAll() ([]*models.SystemAliasPermissionPo, error) {
	var systemAliasPermissionPos []*models.SystemAliasPermissionPo
	err := myself.RepositoryBase.SelectAll(&systemAliasPermissionPos)

	return systemAliasPermissionPos, err
}
