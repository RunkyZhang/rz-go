package repositories

import (
	"time"
)

type messageRepositoryBase struct {
	repositoryBase
}

func (myself *messageRepositoryBase) UpdateById(id int, states string, finished bool, errorMessages string, date time.Time) (int64, error) {
	database, err := myself.getDatabase(nil)
	if nil != err {
		return 0, err
	}
	tableName := myself.getTableName(date)

	keyValues := map[string]interface{}{}
	keyValues["states"] = states
	keyValues["finished"] = finished
	keyValues["errorMessages"] = errorMessages
	keyValues["updatedTime"] = time.Now()
	database = database.Table(tableName).Where("id=?", id).Updates(keyValues)
	if nil != database.Error {
		return 0, database.Error
	}

	return database.RowsAffected, nil
}

func (myself *messageRepositoryBase) SelectByExpireTimeAndFinished(models interface{}, date time.Time) (error) {
	database, err := myself.getDatabase(nil)
	if nil != err {
		return err
	}
	tableName := myself.getTableName(date)

	return database.Table(tableName).Where("finished=0 and deleted=0 and expireTime<? ", time.Now()).Find(models).Error
}
