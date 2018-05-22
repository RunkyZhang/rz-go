package data

import (
	"sync"

	"github.com/garyburd/redigo/redis"
	"time"
	"errors"
)

var (
	redisClient *RedisClient = nil
	lock        sync.Mutex
)

type doFunc func(redis.Conn) (interface{}, error)

type RedisClient struct {
	redisPool *redis.Pool

	RedisClientSettings RedisClientSettings
}

type RedisClientSettings struct {
	PoolMaxActive   int
	PoolMaxIdle     int
	PoolWait        bool
	PoolIdleTimeout time.Duration
	DatabaseId      int
	Password        string
	ConnectTimeout  time.Duration
	Address         string
}

func NewRedisClient(redisClientSettings *RedisClientSettings) (*RedisClient) {
	if nil != redisClient {
		return redisClient
	}

	lock.Lock()
	defer lock.Unlock()

	if nil == redisClient {
		if nil == redisClientSettings {
			redisClientSettings = &RedisClientSettings{
				PoolMaxActive:   10,
				PoolMaxIdle:     1,
				PoolWait:        true,
				PoolIdleTimeout: 180 * time.Second,
				DatabaseId:      0,
				ConnectTimeout:  2000 * time.Second,
			}
		}

		redisClient = &RedisClient{
			RedisClientSettings: *redisClientSettings,
		}

		redisClient.redisPool = &redis.Pool{
			MaxActive:   redisClientSettings.PoolMaxActive,
			MaxIdle:     redisClientSettings.PoolMaxIdle,
			Wait:        redisClientSettings.PoolWait,
			IdleTimeout: redisClientSettings.PoolIdleTimeout,

			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial(
					"tcp",
					redisClientSettings.Address,
					redis.DialDatabase(redisClientSettings.DatabaseId),
					redis.DialPassword(redisClientSettings.Password),
					redis.DialConnectTimeout(redisClientSettings.ConnectTimeout))
				if nil != err {
					return nil, err
				}

				return conn, nil
			},
		}
	}

	return redisClient
}

func (redisClient *RedisClient) StringSet(key string, value string) (error) {
	_, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("SET", key, value)
	})

	return err
}

func (redisClient *RedisClient) StringGet(key string) (string, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("GET", key)
	})

	return redisClient.resultToString(result, err)
}

func (redisClient *RedisClient) safeDo(doFunc doFunc) (interface{}, error) {
	conn := redisClient.redisPool.Get()
	defer conn.Close()

	return doFunc(conn)
}

func (*RedisClient) resultToString(result interface{}, err error) (string, error) {
	if nil != err {
		return "", err
	}

	value, ok := result.(string)
	if !ok {
		return "", errors.New("cannot convert to string")
	}

	return value, nil
}
