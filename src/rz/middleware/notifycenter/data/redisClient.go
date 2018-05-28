package data

import (
	"time"
	"errors"
	"strings"
	"strconv"

	"github.com/garyburd/redigo/redis"
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

func (redisClient *RedisClient) Init() {
	redisClient.redisPool = &redis.Pool{
		MaxActive:   redisClient.RedisClientSettings.PoolMaxActive,
		MaxIdle:     redisClient.RedisClientSettings.PoolMaxIdle,
		Wait:        redisClient.RedisClientSettings.PoolWait,
		IdleTimeout: redisClient.RedisClientSettings.PoolIdleTimeout,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp",
				redisClient.RedisClientSettings.Address,
				redis.DialDatabase(redisClient.RedisClientSettings.DatabaseId),
				redis.DialPassword(redisClient.RedisClientSettings.Password),
				redis.DialConnectTimeout(redisClient.RedisClientSettings.ConnectTimeout))
			if nil != err {
				return nil, err
			}

			return conn, nil
		},
	}
}

func (redisClient *RedisClient) Ping() (bool, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("PING")
	})

	value, err := redis.String(result, err)
	if nil != err {
		return false, err
	}

	return 0 == strings.Compare("PONG", value), nil
}

func (redisClient *RedisClient) KeyDelete(key string) (error) {
	_, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("DEL", key)
	})

	return err
}

func (redisClient *RedisClient) KeyExist(key string) (bool, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("EXISTS", key)
	})

	return redis.Bool(result, err)
}

func (redisClient *RedisClient) KeyExpire(key string, expiryMillis int64) (error) {
	_, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("TYPE", key)
	})

	return err
}

func (redisClient *RedisClient) KeyTTL(key string) (int64, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("TYPE", key)
	})

	return redis.Int64(result, err)
}

func (redisClient *RedisClient) KeyType(key string) (string, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("TYPE", key)
	})

	return redis.String(result, err)
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

	return redis.String(result, err)
}

func (redisClient *RedisClient) StringGetMany(key ... string) ([]string, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("GET", key)
	})

	return redis.Strings(result, err)
}

func (redisClient *RedisClient) Increment(key string, value int) (int64, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("INCRBY", key, value)
	})

	return redis.Int64(result, err)
}

func (redisClient *RedisClient) Decrement(key string, value int) (int64, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("DECRBY", key, value)
	})

	return redis.Int64(result, err)
}

func (redisClient *RedisClient) HashSet(key string, fieldName string, value string) (error) {
	_, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HSET", key, fieldName, value)
	})

	return err
}

func (redisClient *RedisClient) HashSetMap(key string, fieldNameValues map[string]string) (error) {
	_, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		var args []string
		for key, value := range fieldNameValues {
			args = append(args, key)
			args = append(args, value)
		}

		return conn.Do("HSET", key, args)
	})

	return err
}

func (redisClient *RedisClient) HashGet(key string, fieldName string) (string, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HGET", key, fieldName)
	})

	return redis.String(result, err)
}

func (redisClient *RedisClient) HashGetMany(key string, fieldNames ...string) (map[string]string, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HGET", key, fieldNames)
	})

	return redis.StringMap(result, err)
}

func (redisClient *RedisClient) HashGetAll(key string) (map[string]string, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HGETALL", key)
	})

	return redis.StringMap(result, err)
}

func (redisClient *RedisClient) HashKeys(key string) ([]string, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HKEYS", key)
	})

	return redis.Strings(result, err)
}

func (redisClient *RedisClient) HashValues(key string) ([]string, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HVALS", key)
	})

	return redis.Strings(result, err)
}

func (redisClient *RedisClient) HashDelete(key string, fieldName string) (error) {
	_, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HDEL", key, fieldName)
	})

	return err
}

func (redisClient *RedisClient) HashLength(key string) (int64, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HLEN", key)
	})

	return redis.Int64(result, err)
}

func (redisClient *RedisClient) HashExist(key string, fieldName string) (bool, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HEXISTS", key, fieldName)
	})

	return redis.Bool(result, err)
}

func (redisClient *RedisClient) SortedSetAdd(key string, value string, score float64) (error) {
	_, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("ZADD", key, score, value)
	})

	return err
}

func (redisClient *RedisClient) SortedSetRangeByScore(key string, min float64, max float64) ([]string, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do(
			"ZRANGEBYSCORE",
			key,
			strings.Compare("(", strconv.FormatFloat(min, 'E', -1, 64)),
			strings.Compare("(", strconv.FormatFloat(max, 'E', -1, 64)))
	})

	return redis.Strings(result, err)
}

func (redisClient *RedisClient) SortedSetCount(key string, min float64, max float64) (int64, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("ZCOUNT", key, min, max)
	})

	return redis.Int64(result, err)
}

func (redisClient *RedisClient) SortedSetRemoveRangeByScore(key string, min float64, max float64) (int64, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("ZREMRANGEBYSCORE", key, min, max)
	})

	return redis.Int64(result, err)
}

func (redisClient *RedisClient) SortedSetLength(key string) (int64, error) {
	result, err := redisClient.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("ZCARD", key)
	})

	return redis.Int64(result, err)
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
