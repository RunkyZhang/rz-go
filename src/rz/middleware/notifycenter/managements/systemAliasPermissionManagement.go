package managements

import (
	"time"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/repositories"
	"rz/core/common"
	"rz/middleware/notifycenter/exceptions"
)

var (
	SystemAliasPermissionManagement = systemAliasPermissionManagement{
		refreshDuration: 10,
	}

	systemAliasPermissionPos map[string]*models.SystemAliasPermissionPo
)

type systemAliasPermissionManagement struct {
	managementBase

	lastRefreshTime int64
	refreshDuration int
}

func (myself *systemAliasPermissionManagement) Add(systemAliasPermissionPo *models.SystemAliasPermissionPo) (error) {
	err := common.Assert.IsNotNilToError(systemAliasPermissionPo, "systemAliasPermissionPo")
	if nil != err {
		return err
	}
	err = common.Assert.IsNotNilToError(systemAliasPermissionPo.PoBase, "systemAliasPermissionPo.PoBase")
	if nil != err {
		return err
	}

	myself.setPoBase(&systemAliasPermissionPo.PoBase)

	return repositories.SystemAliasPermissionRepository.Insert(systemAliasPermissionPo)
}

func (myself *systemAliasPermissionManagement) Modify(id string, smsPermission *int, mailPermission *int, smsDayFrequency *int, smsHourFrequency *int, smsMinuteFrequency *int) (int64, error) {
	return repositories.SystemAliasPermissionRepository.Update(id, smsPermission, mailPermission, smsDayFrequency, smsHourFrequency, smsMinuteFrequency)
}

func (myself *systemAliasPermissionManagement) GetById(id string) (*models.SystemAliasPermissionPo, error) {
	systemAliasPermissionPos, err := myself.GetAll()
	if nil != err {
		return nil, err
	}

	systemAliasPermissionPo, ok := systemAliasPermissionPos[id]
	if !ok {
		return nil, exceptions.SystemAliasNotExist().AttachError(err).AttachMessage(id)
	}

	return systemAliasPermissionPo, nil
}

func (myself *systemAliasPermissionManagement) GetAll() (map[string]*models.SystemAliasPermissionPo, error) {
	var err error
	if nil == systemAliasPermissionPos {
		systemAliasPermissionPos, err = myself.getAll()
		return systemAliasPermissionPos, err
	}

	if int64(myself.refreshDuration) <= time.Now().Unix()-myself.lastRefreshTime {
		go func() {
			systemAliasPermissionPos, err = myself.getAll()
			if nil == err {
				myself.lastRefreshTime = time.Now().Unix()
			} else {
				common.GetLogging().Error(err, "Failed to get all [SystemAliasPermissionPo]")
			}
		}()
	}

	return systemAliasPermissionPos, nil
}

func (*systemAliasPermissionManagement) getAll() (map[string]*models.SystemAliasPermissionPo, error) {
	systemAliasPermissionPos, err := repositories.SystemAliasPermissionRepository.SelectAll()
	if nil != err {
		return nil, err
	}

	var keyValues = make(map[string]*models.SystemAliasPermissionPo)
	for _, systemAliasPermissionPo := range systemAliasPermissionPos {
		keyValues[systemAliasPermissionPo.SystemAlias] = systemAliasPermissionPo
	}

	return keyValues, nil
}
