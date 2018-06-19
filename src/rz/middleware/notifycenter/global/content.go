package global

import (
	"time"

	"rz/middleware/notifycenter/common"
)

var (
	AsyncWorker = common.NewAsyncJobWorker(5, 1*time.Second)
	WebService  = common.NewWebService(Config.Web.Listen)
	RedisClient = buildRedisClient()
)

func buildRedisClient() (*common.RedisClient) {
	redisClientSettings := &common.RedisClientSettings{
		PoolMaxActive:   10,
		PoolMaxIdle:     1,
		PoolWait:        true,
		PoolIdleTimeout: 180 * time.Second,
		DatabaseId:      Config.Redis.DatabaseId,
		ConnectTimeout:  2000 * time.Second,
		Address:         Config.Redis.Address,
		Password:        Config.Redis.Password,
	}

	redisClient, err := common.NewRedisClient(redisClientSettings)
	common.Assert.IsNilErrorToPanic(err, "[redisClientSettings] is nil")

	return redisClient
}
