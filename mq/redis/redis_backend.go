package redis

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/johnnyeven/libtools/mq"
)

func NewRedisBackend(pool *redis.Pool, prefix string) *RedisBackend {
	return &RedisBackend{
		feedback:     mq.Prefix("feedback", prefix),
		cancellation: mq.Prefix("cancellation", prefix),
		pool:         pool,
	}
}

type RedisBackend struct {
	feedback     string
	cancellation string
	pool         *redis.Pool
}

var _ mq.Backend = (*RedisBackend)(nil)

var ChannelNamePrefix = "mq-channel-"

func (backend *RedisBackend) ListChannel() ([]string, error) {
	list, err := redis.Strings(ConnRedis(backend.pool).Do("KEYS", ChannelNamePrefix+"*"))
	if err != nil {
		return nil, err
	}

	channelList := make([]string, len(list))

	for i, k := range list {
		channelList[i] = strings.TrimPrefix(k, ChannelNamePrefix)
	}

	return channelList, nil
}

func (backend *RedisBackend) ListSubject(channel string) ([]string, error) {
	return redis.Strings(ConnRedis(backend.pool).Do("SMEMBERS", ChannelNamePrefix+channel))
}

func (backend *RedisBackend) RegisterChannel(channel string, subjects []string) error {
	channel = ChannelNamePrefix + channel

	args := []interface{}{channel}
	for _, subject := range subjects {
		args = append(args, subject)
	}

	pingRedis := func() error {
		conn := ConnRedis(backend.pool)

		conn.Send("DEL", channel)
		conn.Send("SADD", args...)
		conn.Send("EXPIRE", channel, 70)

		_, err := conn.Exec()
		return err
	}

	if err := pingRedis(); err != nil {
		return err
	}

	go func() {
		for {
			time.Sleep(1 * time.Minute)
			pingRedis()
		}
	}()

	return nil
}

func (backend *RedisBackend) IsCancelled(id string) (bool, error) {
	i, err := redis.Int(ConnRedis(backend.pool).Do("HEXISTS", backend.cancellation, id))
	return i == 1, err
}

func (backend *RedisBackend) Cancel(id string) error {
	_, err := ConnRedis(backend.pool).Do("HSET", backend.cancellation, id, "1")
	return err
}

func (backend *RedisBackend) ClearCancellation(id string) error {
	_, err := ConnRedis(backend.pool).Do("HDEL", backend.cancellation, id)
	return err
}

func (backend *RedisBackend) FeedBack(taskStatus *mq.TaskStatus) error {
	data, err := json.Marshal(taskStatus)
	if err != nil {
		return err
	}
	_, err = ConnRedis(backend.pool).Do("RPUSH", backend.feedback, data)
	return err
}

func (backend *RedisBackend) GetFeedback() (*mq.TaskStatus, error) {
	ret, err := ConnRedis(backend.pool).Do("BLPOP", backend.feedback, "1")
	if err != nil {
		return nil, err
	}

	if ret == nil {
		return nil, fmt.Errorf("null message received from redis")
	}

	msgPair := ret.([]interface{})

	if string(msgPair[0].([]byte)) != backend.feedback {
		return nil, fmt.Errorf("not a backend message: %v", msgPair[0])
	}

	task := (*mq.TaskStatus)(nil)
	if err := json.Unmarshal(msgPair[1].([]byte), &task); err != nil {
		return nil, err
	}
	return task, nil
}
