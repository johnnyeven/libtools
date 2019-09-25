package task

import (
	"github.com/johnnyeven/libtools/conf"
	"github.com/johnnyeven/libtools/task/kafka"
	"os"
	"strings"
	"fmt"
	"github.com/johnnyeven/libtools/env"
	"github.com/johnnyeven/libtools/task/constants"
	"sync"
	"github.com/johnnyeven/libtools/courier"
	"reflect"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
	"github.com/google/uuid"
	"github.com/johnnyeven/libtools/task/gearman"
	"gopkg.in/robfig/cron.v2"
)

const (
	DefaultCentralChannel  = "service-task-manager.dev"
	DefaultCentralSubject  = "CreateOrUpdateCronTable"
	DefaultFeedbackSubject = "TaskStatusFeedback"
)

type Agent struct {
	ConnectionInfo constants.ConnectionInfo
	Channel        string
	BrokerType     constants.BrokerType `conf:"env"`
	workers        []*Worker
	client         Producer
	backend        *Backend
	jobs           sync.Map

	CentralChannel  string `conf:"env"`
	CentralSubject  string `conf:"env"`
}

func (a *Agent) Init() {
	if a.BrokerType == constants.BROKER_TYPE__GEARMAN {
		a.client = gearman.NewGearmanProducer(a.ConnectionInfo)
	} else if a.BrokerType == constants.BROKER_TYPE__KAFKA {
		a.client = kafka.NewKafkaProducer(a.ConnectionInfo)
	}
}

func (Agent) MarshalDefaults(v interface{}) {
	if agent, ok := v.(*Agent); ok {
		if agent.BrokerType == constants.BROKER_TYPE_UNKNOWN {
			agent.BrokerType = constants.BROKER_TYPE__GEARMAN
		}
		if agent.Channel == "" {
			agent.Channel = getDefaultChannel()
		}
		if agent.ConnectionInfo.Port == 0 {
			agent.ConnectionInfo.Port = 4730
		}
		if agent.ConnectionInfo.Protocol == "" {
			agent.ConnectionInfo.Protocol = "tcp"
		}
		if agent.CentralChannel == "" {
			agent.CentralChannel = DefaultCentralChannel
		}
		if agent.CentralSubject == "" {
			agent.CentralSubject = DefaultCentralSubject
		}
	}
}

func (a *Agent) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Protocol":        "tcp",
		"Host":            conf.RancherInternal("dep-tools", "gearmand"),
		"Port":            4730,
		"BrokerType":      1,
		"CentralChannel":  DefaultCentralChannel,
		"CentralSubject":  DefaultCentralSubject,
	}
}

func getDefaultChannel() string {
	serviceName := "anonymous"
	if projectName, exists := os.LookupEnv("PROJECT_NAME"); exists {
		serviceName = projectName
	}
	return strings.ToLower(fmt.Sprintf("%s.%s", serviceName, env.GetRuntimeEnv()))
}

func (a *Agent) Register(subject string, processor constants.TaskProcessor) {
	a.jobs.Store(subject, processor)
}

func (a *Agent) RegisterRoutes(routes ...*courier.Route) {
	for _, route := range routes {
		if len(route.Operators) == 0 {
			continue
		}

		operatorMetas := courier.ToOperatorMetaList(route.Operators...)

		lastOpIndex := len(operatorMetas) - 1
		lastOpMeta := operatorMetas[lastOpIndex]
		subject := lastOpMeta.Type.Name()

		if cronDescriber, ok := lastOpMeta.Operator.(CronDescriber); ok {
			spec := cronDescriber.CronSpec()
			_, err := cron.Parse(spec)
			if err != nil {
				panic(err)
			}
			if err = a.registerCron(subject, spec); err != nil {
				panic(err)
			}
		}

		a.Register(subject, func(task *constants.Task) (interface{}, error) {
			ctx := context.Background()

			for i, opMeta := range operatorMetas {
				op := reflect.New(opMeta.Type).Interface().(courier.IOperator)

				if err := constants.UnmarshalData(task.Data, op); err != nil {
					return nil, err
				}

				ret, err := op.Output(ctx)
				if err != nil {
					return nil, err
				}

				if i != lastOpIndex {
					if ctxProvider, ok := op.(courier.IContextProvider); ok {
						ctx = context.WithValue(ctx, ctxProvider.ContextKey(), ret)
					}
					continue
				}

				return ret, nil
			}

			return nil, nil
		})
	}
}

func (a *Agent) registerCron(subject string, spec string) error {
	cronTable := &constants.CronTableInfo{}
	cronTable.CronTableID = fmt.Sprintf("%s-%s", a.Channel, subject)
	cronTable.Subject = subject
	cronTable.Channel = a.Channel
	cronTable.Spec = spec
	cronTable.Desc = cronTable.CronTableID

	bytes, err := json.Marshal(cronTable)
	if err != nil {
		return err
	}
	_, err = a.Publish(a.CentralChannel, a.CentralSubject, bytes)
	return err
}

func (a *Agent) SendTask(task *constants.Task) error {
	return a.client.SendTask(task)
}

func (a *Agent) Publish(channel string, subject string, data []byte) (*constants.Task, error) {
	if a.backend == nil {
		return nil, fmt.Errorf("backend not init")
	}
	if a.client == nil {
		return nil, fmt.Errorf("client not init")
	}
	task := &constants.Task{
		ID:         uuid.New().String(),
		Channel:    channel,
		Subject:    subject,
		Data:       data,
		CreateTime: time.Now(),
	}

	err := a.backend.Feedback(task.Init())
	if err != nil {
		return nil, err
	}

	err = a.client.SendTask(task)
	if err != nil {
		return nil, err
	}

	err = a.backend.Feedback(task.Pending())
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (a *Agent) Start(numWorker int) {
	if a.Channel == "" {
		logrus.Panic("channel must be set")
	}

	if a.backend == nil {
		a.backend = NewBackend(a)
	}

	if a.workers == nil {
		a.workers = make([]*Worker, 0)
	}
	for i := 0; i < numWorker; i++ {
		w := NewWorker(a.BrokerType, a.ConnectionInfo)
		a.workers = append(a.workers, w)
		w.Start(a.Channel, a.workerProcessor)
	}
}

func (a *Agent) Stop() {
	if a.workers != nil {
		for _, w := range a.workers {
			w.Stop()
		}
	}
}

func (a *Agent) workerProcessor(task *constants.Task) (interface{}, error) {
	logrus.Debugf("receive task: id=%s, channel=%s, subject=%s", task.ID, task.Channel, task.Subject)
	subject := task.Subject
	p, ok := a.jobs.Load(subject)
	if !ok {
		err := fmt.Errorf("subject %s not registered", subject)
		return nil, err
	}

	if task.Subject == DefaultFeedbackSubject {
		ret, err := p.(constants.TaskProcessor)(task)
		if err != nil {
			return nil, err
		}
		return ret, nil
	} else {
		a.backend.Feedback(task.Processing())
		ret, err := p.(constants.TaskProcessor)(task)
		if err != nil {
			a.backend.Feedback(task.Fail(err))
			return nil, err
		}
		a.backend.Feedback(task.Success(ret))
		return ret, nil
	}
}
