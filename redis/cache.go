package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

// Cache is Redis cache adapter.
type RedisCache struct {
	Pool *redis.Pool // redis connection pool
}

// NewRedisCache create new redis cache with default collection name.
func NewRedisCache() *RedisCache {
	return &RedisCache{}
}

func RedisDo(conn redis.Conn) func(cmd string, args ...interface{}) (interface{}, error) {
	return func(cmd string, args ...interface{}) (interface{}, error) {
		res, err := conn.Do(cmd, args...)
		if err != nil {
			logrus.WithField("redis_cmd", cmd).Warningf("%v", err.Error())
		}
		return res, err
	}
}

func formatMs(dur time.Duration) int64 {
	if dur > 0 && dur < time.Millisecond {
		return 0
	}
	return int64(dur / time.Millisecond)
}

func formatSec(dur time.Duration) int64 {
	if dur > 0 && dur < time.Second {
		return 0
	}
	return int64(dur / time.Second)
}

func usePrecise(dur time.Duration) bool {
	return dur < time.Second || dur%time.Second != 0
}

// Start start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info","dbNum":"0"}
// the cache item in redis are stored forever,
// so no gc operation.
func (cache *RedisCache) Start(config *Redis) error {
	if err := cache.connectInit(config); err != nil {
		return err
	}

	c := cache.Pool.Get()
	defer c.Close()

	return c.Err()
}

// connect to redis.
func (cache *RedisCache) connectInit(config *Redis) error {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.DialTimeout(
			config.Protocol,
			fmt.Sprintf("%s:%d", config.Host, config.Port),
			config.ConnectTimeout,
			config.ReadTimeout,
			config.WriteTimeout,
		)
		if err != nil {
			return
		}

		if config.Password != "" {
			if _, err := RedisDo(c)("AUTH", config.Password.String()); err != nil {
				c.Close()
				return c, err
			}
		}

		_, selectErr := RedisDo(c)("SELECT", config.DB)
		if selectErr != nil {
			c.Close()
			return nil, selectErr
		}

		return
	}

	// initialize a new pool
	cache.Pool = &redis.Pool{
		Dial:        dialFunc,
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: config.IdleTimeout,
		Wait:        true,
	}

	return nil
}

// Redis `SET key value [expiration] NX` command.
// Zero expiration means the key has no expiration time.
func (cache *RedisCache) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	var err error
	if expiration == 0 {
		// Use old `SETNX` to support old Redis versions.
		_, err = redis.String(RedisDo(conn)("setnx", key, value))
	} else {
		if usePrecise(expiration) {
			_, err = redis.String(RedisDo(conn)("set", key, value, "nx", "px", formatMs(expiration)))
		} else {
			_, err = redis.String(RedisDo(conn)("set", key, value, "nx", "ex", formatSec(expiration)))
		}
	}

	if err != nil {
		if err != redis.ErrNil {
			return false, err
		} else {
			return false, nil
		}
	} else {
		return true, err
	}
}

func (cache *RedisCache) Get(key string) ([]byte, error) {
	conn := cache.Pool.Get()
	defer conn.Close()

	return redis.Bytes(RedisDo(conn)("GET", key))
}

func (cache *RedisCache) Incr(key string) (int64, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return redis.Int64(RedisDo(conn)("INCR", key))
}

func (cache *RedisCache) MGet(key []interface{}) (interface{}, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return RedisDo(conn)("MGET", key...)
}

func (cache *RedisCache) MGetValue(keys []interface{}) ([]interface{}, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return redis.Values(RedisDo(conn)("MGET", keys...))
}

func (cache *RedisCache) HSet(key, field string, value interface{}) (interface{}, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return RedisDo(conn)("HSET", key, field, value)
}

func (cache *RedisCache) HSetWithExpire(key string, timeout time.Duration) (interface{}, error) {
	conn := cache.Pool.Get()
	defer conn.Close()

	reply, err := RedisDo(conn)("HSET", key, nil, nil)
	if err != nil {
		return reply, err
	}

	return RedisDo(conn)("EXPIRE", key, formatSec(timeout))
}

func (cache *RedisCache) HMset(value []interface{}) (interface{}, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return RedisDo(conn)("HMSET", value...)
}

func (cache *RedisCache) HGet(key, field string) (interface{}, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return RedisDo(conn)("HGET", key, field)
}

func (cache *RedisCache) HIncrBy(key, field string, value interface{}) (int64, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return redis.Int64(RedisDo(conn)("HINCRBY", key, field, value))
}

func (cache *RedisCache) HMGet(key string, fields []string) (interface{}, error) {
	conn := cache.Pool.Get()
	defer conn.Close()

	var args []interface{}
	args = append(args, key)
	for _, field := range fields {
		args = append(args, field)
	}

	return RedisDo(conn)("HMGET", args...)
}

func (cache *RedisCache) GetString(key string) (string, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return redis.String(RedisDo(conn)("GET", key))
}

func (cache *RedisCache) GetStringMap(key string) (map[string]string, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return redis.StringMap(RedisDo(conn)("HGETALL", key))
}

func (cache *RedisCache) HGetAll(key string) (interface{}, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return RedisDo(conn)("HGETALL", key)
}

func (cache *RedisCache) GetInts(key string) ([]int, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return redis.Ints(RedisDo(conn)("GET", key))
}

func (cache *RedisCache) Expire(key string, timeout time.Duration) error {
	conn := cache.Pool.Get()
	defer conn.Close()
	_, err := RedisDo(conn)("EXPIRE", key, formatSec(timeout))
	return err
}

func (cache *RedisCache) Set(key string, bytes interface{}, timeout time.Duration) error {
	conn := cache.Pool.Get()
	defer conn.Close()
	if timeout == -1 {
		_, err := RedisDo(conn)("SET", key, bytes)
		return err
	} else {
		_, err := RedisDo(conn)("SET", key, bytes, "EX", formatSec(timeout))
		return err
	}
}

func (cache *RedisCache) Del(key string) error {
	conn := cache.Pool.Get()
	defer conn.Close()
	_, err := RedisDo(conn)("DEL", key)

	return err
}

func (cache *RedisCache) Exists(key string) (bool, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	exists, err := redis.Bool(RedisDo(conn)("EXISTS", key))
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func (cache *RedisCache) ZRange(key string, start, end int, withscores bool) ([]string, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	var res []string
	var err error
	if withscores {
		res, err = redis.Strings(RedisDo(conn)("ZRANGE", key, start, end, "withscores"))
	} else {
		res, err = redis.Strings(RedisDo(conn)("ZRANGE", key, start, end))
	}
	return res, err
}

func (cache *RedisCache) ZRangeInts(key string, start, end int, withscores bool) ([]int, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	var res []int
	var err error
	if withscores {
		res, err = redis.Ints(RedisDo(conn)("ZRANGE", key, start, end, "withscores"))
	} else {
		res, err = redis.Ints(RedisDo(conn)("ZRANGE", key, start, end))
	}
	return res, err
}

func (cache *RedisCache) ZRevrange(key string, start, end int, withscores bool) ([]int, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	var res []int
	var err error
	if withscores {
		res, err = redis.Ints(RedisDo(conn)("ZREVRANGE", key, start, end, "withscores"))
	} else {
		res, err = redis.Ints(RedisDo(conn)("ZREVRANGE", key, start, end))
	}

	return res, err
}

func (cache *RedisCache) ZRevrangeByScore(key string, max_num, min_num int, withscores bool, offset, count int) ([]int, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	var res []int
	var err error
	if !withscores {
		res, err = redis.Ints(RedisDo(conn)("ZREVRANGEBYSCORE", key, max_num, min_num, "limit", offset, count))
	} else {
		res, err = redis.Ints(RedisDo(conn)("ZREVRANGEBYSCORE", key, max_num, min_num, "withscores", "limit", offset, count))
	}
	return res, err
}
func (cache *RedisCache) ZRangeByScore(key string, min_num, max_num int64, withscores bool, offset, count int) ([]int, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	var res []int
	var err error
	if withscores {
		res, err = redis.Ints(RedisDo(conn)("ZREVRANGEBYSCORE", key, max_num, min_num, "limit", offset, count))
	} else {
		res, err = redis.Ints(RedisDo(conn)("ZREVRANGEBYSCORE", key, max_num, min_num, "withscores", "limit", offset, count))
	}
	return res, err
}

func (cache *RedisCache) ZScore(key, member string) (interface{}, error) {
	conn := cache.Pool.Get()
	defer conn.Close()

	res, err := RedisDo(conn)("ZSCORE", key, member)
	return res, err
}

func (cache *RedisCache) ZAdd(key string, value, member interface{}) (interface{}, error) {
	conn := cache.Pool.Get()
	defer conn.Close()

	res, err := RedisDo(conn)("ZADD", key, value, member)
	return res, err
}

func (cache *RedisCache) SAdd(key string, items string) (int, error) {
	//var err error
	conn := cache.Pool.Get()
	defer conn.Close()
	res, err := redis.Int(RedisDo(conn)("SADD", key, items))
	return res, err
}

func (cache *RedisCache) SIsMember(key string, items string) (int, error) {
	//var err error
	conn := cache.Pool.Get()
	defer conn.Close()
	res, err := redis.Int(RedisDo(conn)("SISMEMBER", key, items))
	return res, err
}

func (cache *RedisCache) RPush(key string, value interface{}) error {
	conn := cache.Pool.Get()
	defer conn.Close()
	_, err := RedisDo(conn)("RPUSH", key, value)
	return err
}

func (cache *RedisCache) RPushBatch(keys []interface{}) error {
	conn := cache.Pool.Get()
	defer conn.Close()
	_, err := RedisDo(conn)("RPUSH", keys...)
	return err
}

func (cache *RedisCache) LRange(key string, start, end int) ([]interface{}, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	result, err := redis.Values(RedisDo(conn)("LRANGE", key, start, end))
	return result, err
}

func (cache *RedisCache) LRem(key string, value interface{}) error {
	conn := cache.Pool.Get()
	defer conn.Close()
	_, err := RedisDo(conn)("LREM", key, 0, value)
	return err
}

func (cache *RedisCache) TTL(key string) (int, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return redis.Int(RedisDo(conn)("ttl", key))
}
