package repositories

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/common"
	"time"
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

func (myself *systemAliasPermissionRepository) Update(id string, smsPermission *int, mailPermission *int, smsDayFrequency *int, smsHourFrequency *int, smsMinuteFrequency *int) (int64, error) {
	err := common.Assert.IsNotBlankToError(id, "id")
	if nil != err {
		return 0, err
	}

	database, err := myself.GetShardDatabase(nil)
	if nil != err {
		return 0, err
	}

	keyValues := map[string]interface{}{}
	if nil != smsPermission {
		keyValues["smsPermission"] = smsPermission
	}
	if nil != mailPermission {
		keyValues["mailPermission"] = mailPermission
	}
	if nil != smsDayFrequency {
		keyValues["smsDayFrequency"] = smsDayFrequency
	}
	if nil != smsHourFrequency {
		keyValues["smsHourFrequency"] = smsHourFrequency
	}
	if nil != smsMinuteFrequency {
		keyValues["smsMinuteFrequency"] = smsMinuteFrequency
	}
	keyValues["updatedTime"] = time.Now()
	database = database.Where("id=?", id).Updates(keyValues)

	return database.RowsAffected, database.Error
}

func (myself *systemAliasPermissionRepository) SelectAll() ([]*models.SystemAliasPermissionPo, error) {
	var systemAliasPermissionPos []*models.SystemAliasPermissionPo
	err := myself.RepositoryBase.SelectAll(&systemAliasPermissionPos)

	return systemAliasPermissionPos, err
}
