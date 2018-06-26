package global

import (
	"time"
	"sync"
	"fmt"

	"rz/middleware/notifycenter/common"
)

var (
	AsyncWorker                     = common.NewAsyncJobWorker(5, 1*time.Second)
	WebService                      = common.NewWebService(GetConfig().Web.Listen)
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
			fmt.Printf("failed to refresh redis\n")

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
