package redis

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/mq"
	"github.com/johnnyeven/libtools/redis"
)

var r = &redis.Redis{
	Name: "mq_test",
	Host: "redis.staging.g7pay.net",
	Port: 36379,
}

func init() {
	r.MarshalDefaults(r)
	r.Init()
}

func TestRedisBroker(t *testing.T) {
	tt := assert.New(t)

	channel := uuid.New().String()

	broker := NewRedisBroker(r.GetCache().Pool, "mq-")

	taskList := make([]*mq.Task, 0)

	for i := 0; i < 10; i++ {
		task := mq.NewTask(channel, "test", []byte(fmt.Sprintf("data %d", i)))
		taskList = append(taskList, task)

		err := broker.SendTask(task)
		tt.NoError(err)
	}

	for i := 0; i < 10; i++ {
		task, err := broker.GetTask(channel)
		tt.NoError(err)
		tt.Equal(taskList[i].ID, task.ID)
	}
}
