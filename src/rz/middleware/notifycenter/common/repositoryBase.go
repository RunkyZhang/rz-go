package common

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"errors"
	"fmt"
)

var (
	Databases map[string]*gorm.DB
)

func InitDatabases(connectionStrings map[string]string) {
	if nil == Databases {
		Databases = make(map[string]*gorm.DB)
		for key, value := range connectionStrings {
			database, err := gorm.Open("mysql", value)
			if nil != err {
				closeDatabase()
				panic(errors.New("failed to open database, error: " + err.Error()))
			}
			database.DB().SetMaxIdleConns(2)
			database.DB().SetMaxOpenConns(10)

			Databases[key] = database
		}
	}
}

func closeDatabase() {
	for _, value := range Databases {
		value.Close()
	}
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
	err := Assert.IsNotNilToError(po, "po")
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
	err := Assert.IsNotNilToError(po, "po")
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
	err := Assert.IsNotNilToError(po, "po")
	if nil != err {
		return err
	}

	database, err := myself.GetShardDatabase(shardParameters...)
	if nil != err {
		return err
	}

	return database.Where("id=? and deleted=0", id).First(po).Error
}

func (myself *RepositoryBase) SelectAll(pos interface{}, shardParameters ...interface{}) (error) {
	err := Assert.IsNotNilToError(pos, "pos")
	if nil != err {
		return err
	}

	database, err := myself.GetShardDatabase(shardParameters...)
	if nil != err {
		return err
	}

	return database.Where("deleted=0").Find(pos).Error
}

func (myself *RepositoryBase) GetShardDatabase(shardParameters ...interface{}) (*gorm.DB, error) {
	var defaultDatabaseKey string
	if nil != myself.GetDatabaseKeyFunc {
		defaultDatabaseKey = myself.GetDatabaseKeyFunc(shardParameters...)
	} else {
		defaultDatabaseKey = myself.DefaultDatabaseKey
	}

	database, ok := Databases[defaultDatabaseKey]
	if !ok {
		return nil, errors.New(fmt.Sprintf("failed to get database(%s)", defaultDatabaseKey))
	}

	var tableName string
	if nil != myself.GetTableNameFunc {
		tableName = myself.GetTableNameFunc(shardParameters...)
	} else {
		tableName = myself.RawTableName
	}

	return database.Table(tableName), nil
}
