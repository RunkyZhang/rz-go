package healths

import (
	"rz/middleware/notifycenter/common"
	"github.com/jinzhu/gorm"
)

func NewMySQLHealthIndicator(databases map[string]*gorm.DB) (*MySQLHealthIndicator, error) {
	err := common.Assert.IsNotNilToError(databases, "databases")
	if nil != err {
		return nil, err
	}

	mysqlHealthIndicator := &MySQLHealthIndicator{
		databases: databases,
	}

	return mysqlHealthIndicator, nil
}

type MySQLHealthIndicator struct {
	databases map[string]*gorm.DB
}

func (myself *MySQLHealthIndicator) Indicate() (*common.HealthReport) {
	healthReport := &common.HealthReport{
		Name:   "Ping MySQL",
		Type:   "MySQL",
		Level:  1,
		Detail: make(map[string]interface{}),
	}

	ok := true
	for key, database := range myself.databases {
		err := database.DB().Ping()
		if nil == err {
			healthReport.Detail[key] = "Ok"
		} else {
			ok = false
			healthReport.Detail[key] = err.Error()
		}

	}
	healthReport.Ok = ok

	return healthReport
}
