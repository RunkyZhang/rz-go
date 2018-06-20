package healths

import "rz/middleware/notifycenter/common"

func NewRedisHealthIndicator(redisClient *common.RedisClient) (*RedisHealthIndicator, error) {
	err := common.Assert.IsNotNilToError(redisClient, "redisClient")
	if nil != err {
		return nil, err
	}

	redisHealthIndicator := &RedisHealthIndicator{
		redisClient: redisClient,
	}

	return redisHealthIndicator, nil
}

type RedisHealthIndicator struct {
	redisClient *common.RedisClient
}

func (myself *RedisHealthIndicator) Indicate() (*common.HealthReport) {
	ok, err := myself.redisClient.Ping()

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
