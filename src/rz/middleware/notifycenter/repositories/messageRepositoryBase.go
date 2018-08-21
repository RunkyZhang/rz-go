package repositories

import (
	"time"
	"fmt"

	"rz/core/common"
)

type MessageRepositoryBase struct {
	repositoryBase
}

func (myself *MessageRepositoryBase) UpdateStatesById(id int64, state string, errorMessage string, providerIds string, finished *bool, finishedTime *time.Time) (int64, error) {
	database, err := myself.GetShardDatabase(id)
	if nil != err {
		return 0, err
	}

	var parameters []interface{}
	setSql := "`states`=CONCAT(`states`,?)"
	parameters = append(parameters, "+"+state)
	setSql += ", `updatedTime`=?"
	parameters = append(parameters, time.Now())
	if "" != errorMessage {
		setSql += ", `errorMessages`=CONCAT(`errorMessages`,?)"
		parameters = append(parameters, errorMessage)
	}
	if "" != providerIds {
		setSql += ", `providerIds`=?"
		parameters = append(parameters, providerIds)
	}
	if nil != finished {
		setSql += ", `finished`=?"
		parameters = append(parameters, finished)
	}
	if nil != finishedTime {
		setSql += ", `finishedTime`=?"
		parameters = append(parameters, finishedTime)
	}
	parameters = append(parameters, id)

	sql := fmt.Sprintf("UPDATE `%s` SET %s WHERE id=?", myself.GetTableNameFunc(id), setSql)
	ravDatabase := database.Exec(sql, parameters...)

	return ravDatabase.RowsAffected, ravDatabase.Error
}

func (myself *MessageRepositoryBase) UpdateDisableById(id int64, disable bool) (int64, error) {
	database, err := myself.GetShardDatabase(id)
	if nil != err {
		return 0, err
	}

	keyValues := map[string]interface{}{}
	keyValues["disable"] = disable
	keyValues["updatedTime"] = time.Now()
	database = database.Where("id=?", id).Updates(keyValues)

	return database.RowsAffected, database.Error
}

func (myself *MessageRepositoryBase) selectByExpireTimeAndFinished(pos interface{}, year int) (error) {
	err := common.Assert.IsTrueToError(nil != pos, "nil != pos")
	if nil != err {
		return err
	}

	database, err := myself.GetShardDatabase(year)
	if nil != err {
		return err
	}

	return database.Where("finished=0 and deleted=0 and expireTime<? ", time.Now()).Find(pos).Error
}

func (myself *MessageRepositoryBase) selectById(id int64, po interface{}) (error) {
	err := common.Assert.IsTrueToError(nil != po, "nil != po")
	if nil != err {
		return err
	}

	return myself.RepositoryBase.SelectById(id, po, id)
}

func (myself *MessageRepositoryBase) selectByIds(ids []int64, pos interface{}, year int) (error) {
	err := common.Assert.IsTrueToError(nil != pos, "nil != pos")
	if nil != err {
		return err
	}

	return myself.RepositoryBase.SelectByIds(ids, pos, year)
}

func (myself *MessageRepositoryBase) getDatabaseKey(shardParameters ...interface{}) (string) {
	return myself.DefaultDatabaseKey
}

func (myself *MessageRepositoryBase) getTableName(shardParameters ...interface{}) (string) {
	if nil == shardParameters || 0 == len(shardParameters) {
		return ""
	}

	id, ok := shardParameters[0].(int64)
	if ok && 999 < id {
		value := common.Int64ToString(id)
		return myself.RawTableName + "_" + value[0:4]
	}

	year, ok := shardParameters[0].(int)
	if ok {
		return myself.RawTableName + "_" + common.Int32ToString(year)
	}

	return ""
}
