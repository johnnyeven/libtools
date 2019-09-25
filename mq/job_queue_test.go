package mq_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/mq"
	mq_redis "github.com/johnnyeven/libtools/mq/redis"
	"github.com/johnnyeven/libtools/redis"
)

var r = &redis.Redis{
	Host: "redis.staging.g7pay.net",
	Port: 36379,
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)

	r.MarshalDefaults(r)
	r.Init()
}

func TestJobQueue(t *testing.T) {
	tt := assert.New(t)

	receiver := mq.NewJobQueue(
		mq_redis.NewRedisBroker(r.GetCache().Pool, "mq-"),
		mq_redis.NewRedisBackend(r.GetCache().Pool, "mq-"),
	)

	for i := 0; i < 10; i++ {
		task, err := receiver.Publish("some-test", "test2", []byte(fmt.Sprintf("data %d", i)))
		tt.NoError(err)
		fmt.Println(task.ID)

		if i%2 == 0 {
			receiver.Cancel(task.ID)
		}
	}

	receiver.RegisterReceiver(func(status *mq.TaskStatus) error {
		fmt.Println(status.ID, status.Status)
		return nil
	})

	receiver.StartReceiver(3)

	worker := mq.NewJobQueue(
		mq_redis.NewRedisBroker(r.GetCache().Pool, "mq-"),
		mq_redis.NewRedisBackend(r.GetCache().Pool, "mq-"),
	)

	worker.Register("test2", func(task *mq.Task) (interface{}, error) {
		time.Sleep(1 * time.Second)
		logrus.Printf("%s %s %s", task.Subject, task.ID, string(task.Data))
		return nil, nil
	})

	worker.StartWorker("some-test", 2)

	time.Sleep(6 * time.Second)
}
