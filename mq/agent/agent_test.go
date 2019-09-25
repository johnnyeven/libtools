package agent

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"

	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/mq"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

type DoSomeThing struct {
	Int int
}

func (op DoSomeThing) Output(ctx context.Context) (interface{}, error) {
	logrus.Println("running: ", op.Int)
	time.Sleep(50 * time.Millisecond)
	if op.Int%2 == 0 {
		return nil, fmt.Errorf("error %d", op.Int)
	}
	return op.Int, nil
}

func TestAgent(t *testing.T) {
	agent := Agent{
		Host: "redis.staging.g7pay.net",
		Port: 36379,
	}
	agent.MarshalDefaults(&agent)
	agent.Init()

	agent.RegisterReceiver(func(status *mq.TaskStatus) error {
		logrus.Infoln(
			"received:",
			status.Channel,
			status.Subject,
			status.ID,
			status.Status,
			string(status.Result),
			string(status.Traceback),
		)
		return nil
	})
	agent.StartReceiver()

	agent.RegisterRoutes(courier.NewRouter(DoSomeThing{}).Route())
	agent.StartWorker()

	for i := 0; i < 10; i++ {
		subject, data, _ := SubjectAndDataFromValue(&DoSomeThing{
			Int: i,
		})
		task, _ := agent.Publish(agent.Name, subject, data)
		logrus.Infoln(task.ID, string(task.Data))
		if i%4 == 0 {
			agent.Cancel(task.ID)
		}
	}

	spew.Dump(agent.ListChannel())

	time.Sleep(5 * time.Second)
}
