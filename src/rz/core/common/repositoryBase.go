package common

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"errors"
	"fmt"
	"sync"
)

var (
	connectionStrings map[string]string
	databases         map[string]*gorm.DB
	databaseLock      sync.Mutex
)

func SetConnectionStrings(keyConnectionStrings map[string]string) {
	connectionStrings = keyConnectionStrings

	GetDatabases()
}

func GetDatabases() (map[string]*gorm.DB) {
	if nil != databases {
		return databases
	}

	databaseLock.Lock()
	defer databaseLock.Unlock()

	if nil != databases {
		return databases
	}

	databases = make(map[string]*gorm.DB)
	for key, value := range connectionStrings {
		database, err := gorm.Open("mysql", value)
		if nil != err {
			CloseDatabase()
			panic(errors.New(fmt.Sprintf("Failed to open database; error: %s", err.Error())))
		}
		database.DB().SetMaxIdleConns(2)
		database.DB().SetMaxOpenConns(10)

		databases[key] = database
	}

	return databases
}

func CloseDatabase() {
	for _, value := range databases {
		value.Close()
	}

	databases = nil
}

type getDatabaseKeyFunc func(...interface{}) (string)
type getTableNameFunc func(...interface{}) (string)

type RepositoryBase struct {
	DefaultDatabaseKey string
	RawTableName       string
	GetDatabaseKeyFunc getDatabaseKeyFunc
	GetTableNameFunc   getTableNameFunc
}

func (myself *RepositoryBase) Insert(po interface{}, shardParameters ...interface{}) (error) {
	err := Assert.IsTrueToError(nil != po, "nil != po")
	if nil != err {
		return err
	}

	database, err := myself.GetShardDatabase(shardParameters...)
	if nil != err {
		return err
	}

	return database.Create(po).Error
}

func (myself *RepositoryBase) Update(po interface{}, shardParameters ...interface{}) (error) {
	err := Assert.IsTrueToError(nil != po, "nil != po")
	if nil != err {
		return err
	}

	database, err := myself.GetShardDatabase(shardParameters...)
	if nil != err {
		return err
	}

	return database.Update(po).Error
}

func (myself *RepositoryBase) SelectById(id interface{}, po interface{}, shardParameters ...interface{}) (error) {
	err := Assert.IsTrueToError(nil != po, "nil != po")
	if nil != err {
		return err
	}

	database, err := myself.GetShardDatabase(shardParameters...)
	if nil != err {
		return err
	}

	return database.Where("id=? AND deleted=0", id).First(po).Error
}

func (myself *RepositoryBase) SelectByIds(ids interface{}, pos interface{}, shardParameters ...interface{}) (error) {
	err := Assert.IsTrueToError(nil != pos, "nil != pos")
	if nil != err {
		return err
	}

	database, err := myself.GetShardDatabase(shardParameters...)
	if nil != err {
		return err
	}

	return database.Where("id IN (?) AND deleted=0", ids).Find(pos).Error
}

func (myself *RepositoryBase) SelectAll(pos interface{}, shardParameters ...interface{}) (error) {
	err := Assert.IsTrueToError(nil != pos, "nil != pos")
	if nil != err {
		return err
	}

	database, err := myself.GetShardDatabase(shardParameters...)
	if nil != err {
		return err
	}

	return database.Where("deleted=0").Find(pos).Error
}

func (myself *RepositoryBase) Count(extend int, shardParameters ...interface{}) (int64, error) {
	database, err := myself.GetShardDatabase(shardParameters...)
	if nil != err {
		return 0, err
	}

	var count int64
	err = database.Where("deleted=0", extend).Count(&count).Error

	return count, err
}

func (myself *RepositoryBase) GetShardDatabase(shardParameters ...interface{}) (*gorm.DB, error) {
	var defaultDatabaseKey string
	if nil != myself.GetDatabaseKeyFunc {
		defaultDatabaseKey = myself.GetDatabaseKeyFunc(shardParameters...)
	} else {
		defaultDatabaseKey = myself.DefaultDatabaseKey
	}

	database, ok := GetDatabases()[defaultDatabaseKey]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Failed to get database(%s)", defaultDatabaseKey))
	}

	var tableName string
	if nil != myself.GetTableNameFunc {
		tableName = myself.GetTableNameFunc(shardParameters...)
	} else {
		tableName = myself.RawTableName
	}

	return database.Table(tableName), nil
}
