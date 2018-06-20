package repositories

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"errors"
	"fmt"
	"rz/middleware/notifycenter/common"
)

var (
	Databases map[string]*gorm.DB
)

func Init(connectionStrings map[string]string) {
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

type repositoryBase struct {
	defaultDatabaseKey string
	rawTableName       string
	getDatabaseKeyFunc getDatabaseKeyFunc
	getTableNameFunc   getTableNameFunc
}

func (myself *repositoryBase) Insert(po interface{}, shardingParameters ...interface{}) (error) {
	err := common.Assert.IsNotNilToError(po, "po")
	if nil != err {
		return err
	}

	database, err := myself.getShardingDatabase(shardingParameters...)
	if nil != err {
		return err
	}

	return database.Create(po).Error
}

func (myself *repositoryBase) Update(po interface{}, shardingParameters ...interface{}) (error) {
	err := common.Assert.IsNotNilToError(po, "po")
	if nil != err {
		return err
	}

	database, err := myself.getShardingDatabase(shardingParameters...)
	if nil != err {
		return err
	}

	return database.Update(po).Error
}

func (myself *repositoryBase) SelectById(id int, po interface{}, shardingParameters ...interface{}) (error) {
	err := common.Assert.IsNotNilToError(po, "po")
	if nil != err {
		return err
	}

	database, err := myself.getShardingDatabase(shardingParameters...)
	if nil != err {
		return err
	}

	return database.Where("id=? and deleted=0", id).First(po).Error
}

func (myself *repositoryBase) SelectAll(pos interface{}, shardingParameters ...interface{}) (error) {
	err := common.Assert.IsNotNilToError(pos, "pos")
	if nil != err {
		return err
	}

	database, err := myself.getShardingDatabase(shardingParameters...)
	if nil != err {
		return err
	}

	return database.Where("deleted=0").Find(pos).Error
}

func (myself *repositoryBase) getShardingDatabase(shardingParameters ...interface{}) (*gorm.DB, error) {
	var defaultDatabaseKey string
	if nil != myself.getDatabaseKeyFunc {
		defaultDatabaseKey = myself.getDatabaseKeyFunc(shardingParameters...)
	} else {
		defaultDatabaseKey = myself.defaultDatabaseKey
	}

	database, ok := Databases[defaultDatabaseKey]
	if !ok {
		return nil, errors.New(fmt.Sprintf("failed to get database(%s)", defaultDatabaseKey))
	}

	var tableName string
	if nil != myself.getTableNameFunc {
		tableName = myself.getTableNameFunc(shardingParameters...)
	} else {
		tableName = myself.rawTableName
	}

	return database.Table(tableName), nil
}
