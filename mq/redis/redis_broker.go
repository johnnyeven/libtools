package redis

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"

	"github.com/johnnyeven/libtools/mq"
)

func NewRedisBroker(pool *redis.Pool, prefix string) *RedisBroker {
	return &RedisBroker{
		prefix: prefix,
		pool:   pool,
	}
}

type RedisBroker struct {
	prefix string
	pool   *redis.Pool
}

var _ mq.Broker = (*RedisBroker)(nil)

func (broker *RedisBroker) SendTask(task *mq.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	_, err = ConnRedis(broker.pool).Do("RPUSH", mq.Prefix(task.Channel, broker.prefix), data)
	return err
}

func (broker *RedisBroker) GetTask(channel string) (*mq.Task, error) {
	channel = mq.Prefix(channel, broker.prefix)

	ret, err := ConnRedis(broker.pool).Do("BLPOP", channel, "1")
	if err != nil {
		return nil, err
	}

	if ret == nil {
		return nil, fmt.Errorf("null message received from redis")
	}

	msgPair := ret.([]interface{})

	if string(msgPair[0].([]byte)) != channel {
		return nil, fmt.Errorf("not a broker message: %v", msgPair[0])
	}

	task := (*mq.Task)(nil)
	if err := json.Unmarshal(msgPair[1].([]byte), &task); err != nil {
		return nil, err
	}
	return task, nil
}
