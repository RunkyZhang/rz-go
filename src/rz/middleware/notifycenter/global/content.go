package global

import (
	"time"
	"sync"

	"rz/middleware/notifycenter/common"
)

var (
	Version                         = "2018-06-28-11:55"
	AsyncWorker                     = common.NewAsyncJobWorker(5, 1*time.Second)
	WebService                      = common.NewWebService(GetConfig().Web.Listen)
	HttpClient                      = common.NewHttpClient(nil)
	redisClient *common.RedisClient = nil
	redisLock   sync.Mutex
)

func GetRedisClient() (*common.RedisClient) {
	if nil != redisClient {
		return redisClient
	}

	redisLock.Lock()
	defer redisLock.Unlock()

	if nil != redisClient {
		return redisClient
	}

	redisClientSettings := common.RedisClientSettings{
		PoolMaxActive:   10,
		PoolMaxIdle:     1,
		PoolWait:        true,
		PoolIdleTimeout: 180 * time.Second,
		DatabaseId:      GetConfig().Redis.DatabaseId,
		ConnectTimeout:  2000 * time.Second,
		Address:         GetConfig().Redis.Address,
		Password:        GetConfig().Redis.Password,
	}
	redisClient = common.NewRedisClient(redisClientSettings)

	return redisClient
}

func RefreshRedis() {
	flag := false
	oldRedisClient := redisClient

	defer func() {
		value := recover()
		if nil != value {
			common.GetLogging().Error(value, "failed to refresh redis")

			if !flag {
				redisClient = oldRedisClient
			}
		}
	}()

	redisClient = nil
	GetRedisClient()
	flag = true
	oldRedisClient.Close()
}

func RefreshDatabases() {
	common.SetConnectionStrings(GetConfig().Databases)
	common.CloseDatabase()
	common.GetDatabases()
}
