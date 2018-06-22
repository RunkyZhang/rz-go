package repositories

import (
	"time"
	"rz/middleware/notifycenter/common"
)

type MessageRepositoryBase struct {
	repositoryBase
}

func (myself *MessageRepositoryBase) UpdateById(id int, states string, finished bool, finishedTime time.Time, errorMessages string, date time.Time) (int64, error) {
	database, err := myself.getShardDatabase(date)
	if nil != err {
		return 0, err
	}

	keyValues := map[string]interface{}{}
	keyValues["states"] = states
	keyValues["finished"] = finished
	keyValues["finishedTime"] = finishedTime
	keyValues["errorMessages"] = errorMessages
	keyValues["updatedTime"] = time.Now()
	database = database.Where("id=?", id).Updates(keyValues)

	return database.RowsAffected, database.Error
}

func (myself *MessageRepositoryBase) SelectByExpireTimeAndFinished(pos interface{}, date time.Time) (error) {
	err := common.Assert.IsNotNilToError(pos, "pos")
	if nil != err {
		return err
	}

	database, err := myself.getShardDatabase(nil)
	if nil != err {
		return err
	}

	return database.Where("finished=0 and deleted=0 and expireTime<? ", time.Now()).Find(pos).Error
}
