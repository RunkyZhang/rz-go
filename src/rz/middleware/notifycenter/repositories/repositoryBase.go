package repositories

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"rz/middleware/notifycenter/common"
	"errors"
	"fmt"
)

var (
	Databases map[string]*gorm.DB
)

func Init(connectionStrings map[string]string) {
	if nil == Databases {
		Databases = make(map[string]*gorm.DB)
		common.Assert.IsNotNil(connectionStrings, "connectionStrings")

		for key, value := range connectionStrings {
			//database, err := gorm.Open("mysql", "ua_notifycenter:ekIxrgWsJ03u@tcp(10.0.34.44:3306)/notifycenter")
			database, err := gorm.Open("mysql", value)
			if nil != err {
				closeDatabase()
				panic("failed to open database, error: " + err.Error())
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
	database, err := myself.getDatabase(shardingParameters...)
	if nil != err {
		return err
	}
	tableName := myself.getTableName(shardingParameters...)

	return database.Table(tableName).Create(po).Error
}

func (myself *repositoryBase) Update(po interface{}, shardingParameters ...interface{}) (error) {
	database, err := myself.getDatabase(shardingParameters...)
	if nil != err {
		return err
	}
	tableName := myself.getTableName(shardingParameters...)

	return database.Table(tableName).Update(po).Error
}

func (myself *repositoryBase) SelectById(id int, po interface{}, shardingParameters ...interface{}) (error) {
	database, err := myself.getDatabase(shardingParameters...)
	if nil != err {
		return err
	}
	tableName := myself.getTableName(shardingParameters...)

	return database.Table(tableName).Where("id=? and deleted=0", id).First(po).Error
}

func (myself *repositoryBase) SelectAll(pos interface{}, shardingParameters ...interface{}) (error) {
	database, err := myself.getDatabase(shardingParameters...)
	if nil != err {
		return err
	}
	tableName := myself.getTableName(shardingParameters...)

	return database.Table(tableName).Where("deleted=0").Find(pos).Error
}

func (myself *repositoryBase) getDatabase(shardingParameters ...interface{}) (*gorm.DB, error) {
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

	return database, nil
}

func (myself *repositoryBase) getTableName(shardingParameters ...interface{}) (string) {
	if nil != myself.getTableNameFunc {
		return myself.getTableNameFunc(shardingParameters...)
	} else {
		return myself.defaultDatabaseKey
	}
}
