package redis

import (
	"github.com/gomodule/redigo/redis"
)

func (cache *RedisCache) HGetString(key, field string) (string, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return redis.String(RedisDo(conn)("HGET", key, field))
}

func (cache *RedisCache) HGetUint64(key, field string) (uint64, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return redis.Uint64(RedisDo(conn)("HGET", key, field))
}

func (cache *RedisCache) HGetInt64(key, field string) (int64, error) {
	conn := cache.Pool.Get()
	defer conn.Close()
	return redis.Int64(RedisDo(conn)("HGET", key, field))
}
