package repositories

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	database *gorm.DB
)

func init() {
	var err error
	database, err = gorm.Open("mysql", "ua_notifycenter:ekIxrgWsJ03u@tcp(10.0.34.44:3306)/notifycenter")
	if nil != err{
		panic("failed to open database, error: " + err.Error())
	}
	database.DB().SetMaxIdleConns(2)
	database.DB().SetMaxOpenConns(10)
}
