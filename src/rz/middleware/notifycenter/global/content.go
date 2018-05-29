package global

import (
	"sync"
	"time"

	"rz/middleware/notifycenter/common"
)

var (
	redisClient *common.RedisClient = nil
	redisLock   sync.Mutex
)

func GetRedisClient() (*common.RedisClient) {
	if nil != redisClient {
		return redisClient
	}

	redisLock.Lock()
	defer redisLock.Unlock()

	if nil == redisClient {
		redisClientSettings := common.RedisClientSettings{
			PoolMaxActive:   10,
			PoolMaxIdle:     1,
			PoolWait:        true,
			PoolIdleTimeout: 180 * time.Second,
			DatabaseId:      Config.Redis.DatabaseId,
			ConnectTimeout:  2000 * time.Second,
			Address:         Config.Redis.Address,
			Password:        Config.Redis.Password,
		}

		redisClient = &common.RedisClient{
			RedisClientSettings: redisClientSettings,
		}
		redisClient.Init()
	}

	return redisClient
}
