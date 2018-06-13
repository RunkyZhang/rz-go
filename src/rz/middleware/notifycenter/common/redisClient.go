package common

import (
	"time"
	"errors"
	"strings"
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

func (myself *RedisClient) Init() {
	myself.redisPool = &redis.Pool{
		MaxActive:   myself.RedisClientSettings.PoolMaxActive,
		MaxIdle:     myself.RedisClientSettings.PoolMaxIdle,
		Wait:        myself.RedisClientSettings.PoolWait,
		IdleTimeout: myself.RedisClientSettings.PoolIdleTimeout,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp",
				myself.RedisClientSettings.Address,
				redis.DialDatabase(myself.RedisClientSettings.DatabaseId),
				redis.DialPassword(myself.RedisClientSettings.Password),
				redis.DialConnectTimeout(myself.RedisClientSettings.ConnectTimeout))
			if nil != err {
				return nil, err
			}

			return conn, nil
		},
	}
}

func (myself *RedisClient) Ping() (bool, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("PING")
	})

	value, err := redis.String(result, err)
	if nil != err {
		return false, err
	}

	return 0 == strings.Compare("PONG", value), nil
}

func (myself *RedisClient) KeyDelete(keys ...string) (error) {
	_, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		var args []interface{}
		if nil != keys {
			for _, key := range keys {
				args = append(args, key)
			}
		}

		return conn.Do("DEL", args...)
	})

	return err
}

func (myself *RedisClient) KeyExist(key string) (bool, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("EXISTS", key)
	})

	return redis.Bool(result, err)
}

func (myself *RedisClient) KeyExpire(key string, expiryMillis int64) (error) {
	_, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("TYPE", key)
	})

	return err
}

func (myself *RedisClient) KeyTTL(key string) (int64, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("TYPE", key)
	})

	return redis.Int64(result, err)
}

func (myself *RedisClient) KeyType(key string) (string, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("TYPE", key)
	})

	return redis.String(result, err)
}

func (myself *RedisClient) StringSet(key string, value string) (error) {
	_, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("SET", key, value)
	})

	return err
}

func (myself *RedisClient) StringGet(key string) (string, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("GET", key)
	})

	return redis.String(result, err)
}

func (myself *RedisClient) StringGetMany(keys ... string) ([]string, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		var args []interface{}
		if nil != keys {
			for _, key := range keys {
				args = append(args, key)
			}
		}

		return conn.Do("GET", args...)
	})

	return redis.Strings(result, err)
}

func (myself *RedisClient) Increment(key string, value int) (int64, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("INCRBY", key, value)
	})

	return redis.Int64(result, err)
}

func (myself *RedisClient) Decrement(key string, value int) (int64, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("DECRBY", key, value)
	})

	return redis.Int64(result, err)
}

func (myself *RedisClient) HashSet(key string, fieldName string, value string) (error) {
	_, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HSET", key, fieldName, value)
	})

	return err
}

func (myself *RedisClient) HashSetMap(key string, fieldNameValues map[string]string) (error) {
	_, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		var args []interface{}
		args = append(args, key)
		for fieldName, value := range fieldNameValues {
			args = append(args, fieldName)
			args = append(args, value)
		}

		return conn.Do("HSET", args...)
	})

	return err
}

func (myself *RedisClient) HashGet(key string, fieldName string) (string, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HGET", key, fieldName)
	})

	return redis.String(result, err)
}

func (myself *RedisClient) HashGetMany(key string, fieldNames ...string) ([]string, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		var args []interface{}
		args = append(args, key)
		if nil != fieldNames {
			for _, fieldName := range fieldNames {
				args = append(args, fieldName)
			}
		}
		return conn.Do("HMGET", args...)
	})

	return redis.Strings(result, err)
}

func (myself *RedisClient) HashGetAll(key string) (map[string]string, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HGETALL", key)
	})

	return redis.StringMap(result, err)
}

func (myself *RedisClient) HashKeys(key string) ([]string, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HKEYS", key)
	})

	return redis.Strings(result, err)
}

func (myself *RedisClient) HashValues(key string) ([]string, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HVALS", key)
	})

	return redis.Strings(result, err)
}

func (myself *RedisClient) HashDelete(key string, fieldNames ...string) (int64, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		var args []interface{}
		args = append(args, key)
		if nil != fieldNames {
			for _, fieldName := range fieldNames {
				args = append(args, fieldName)
			}
		}

		return conn.Do("HDEL", args...)
	})

	return redis.Int64(result, err)
}

func (myself *RedisClient) HashLength(key string) (int64, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HLEN", key)
	})

	return redis.Int64(result, err)
}

func (myself *RedisClient) HashExist(key string, fieldName string) (bool, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HEXISTS", key, fieldName)
	})

	return redis.Bool(result, err)
}

func (myself *RedisClient) SortedSetAdd(key string, value string, score float64) (error) {
	_, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("ZADD", key, score, value)
	})

	return err
}

func (myself *RedisClient) SortedSetRangeByScore(key string, min float64, max float64) ([]string, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("ZRANGEBYSCORE", key, min, max)
	})

	return redis.Strings(result, err)
}

func (myself *RedisClient) SortedSetCount(key string, min float64, max float64) (int64, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("ZCOUNT", key, min, max)
	})

	return redis.Int64(result, err)
}

func (myself *RedisClient) SortedSetRemoveRangeByScore(key string, min float64, max float64) (int64, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("ZREMRANGEBYSCORE", key, min, max)
	})

	return redis.Int64(result, err)
}

func (myself *RedisClient) SortedSetRemoveByValue(key string, values ...string) (int64, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		var args []interface{}
		args = append(args, key)
		if nil != values {
			for _, value := range values {
				args = append(args, value)
			}
		}

		return conn.Do("ZREM", args...)
	})

	return redis.Int64(result, err)
}

func (myself *RedisClient) SortedSetLength(key string) (int64, error) {
	result, err := myself.safeDo(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("ZCARD", key)
	})

	return redis.Int64(result, err)
}

func (myself *RedisClient) safeDo(doFunc doFunc) (interface{}, error) {
	conn := myself.redisPool.Get()
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
