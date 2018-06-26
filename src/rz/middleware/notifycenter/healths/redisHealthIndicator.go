package healths

import (
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/global"
)

type RedisHealthIndicator struct {
}

func (myself *RedisHealthIndicator) Indicate() (*common.HealthReport) {
	ok, err := global.GetRedisClient().Ping()

	var message string
	if nil != err {
		message = err.Error()
	}

	return &common.HealthReport{
		Ok:      ok,
		Name:    "Ping Redis",
		Type:    "Redis",
		Level:   1,
		Message: message,
	}
}
