package healths

import (
	"rz/middleware/notifycenter/common"
)

type MySQLHealthIndicator struct {
}

func (myself *MySQLHealthIndicator) Indicate() (*common.HealthReport) {
	healthReport := &common.HealthReport{
		Name:   "Ping MySQL",
		Type:   "MySQL",
		Level:  1,
		Detail: make(map[string]interface{}),
	}

	ok := true
	for key, database := range common.GetDatabases() {
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
