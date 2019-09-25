package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"gopkg.in/robfig/cron.v2"

	"github.com/johnnyeven/libtools/conf"
	"github.com/johnnyeven/libtools/conf/presets"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/env"
	"github.com/johnnyeven/libtools/mq"
	mq_redis "github.com/johnnyeven/libtools/mq/redis"
	"github.com/johnnyeven/libtools/reflectx"
	"github.com/johnnyeven/libtools/timelib"
)

var (
	InvalidDeferTaskTime = errors.New("延迟执行任务时间错误")
)

type Agent struct {
	Name string

	CentreChannel string `conf:"env"`

	Protocol       string
	Host           string `conf:"env,upstream"`
	Port           int32 `conf:"env"`
	Password       presets.Password `conf:"env"`
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxActive      int
	MaxIdle        int
	Wait           bool
	DB             int

	NumWorkers int `conf:"env"`

	Prefix string

	pool  *redis.Pool
	queue *mq.JobQueue
}

func (Agent) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Host": conf.RancherInternal("db-routers", "task"),
		"Port": 6379,
	}
}

func getDefaultName() string {
	serviceName := "anonymous"
	if projectName, exists := os.LookupEnv("PROJECT_NAME"); exists {
		serviceName = projectName
	}
	return strings.ToLower(fmt.Sprintf("%s.%s", serviceName, env.GetRuntimeEnv()))
}

func (Agent) MarshalDefaults(v interface{}) {
	if agent, ok := v.(*Agent); ok {
		if agent.CentreChannel == "" {
			agent.CentreChannel = "service-task-mgr.dev"
		}
		if agent.Name == "" {
			agent.Name = getDefaultName()
		}
		if agent.Prefix == "" {
			agent.Prefix = "mq-"
		}
		if agent.NumWorkers == 0 {
			agent.NumWorkers = 1
		}
		if agent.Protocol == "" {
			agent.Protocol = "tcp"
		}
		if agent.Port == 0 {
			agent.Port = 6379
		}
		if agent.Password == "" {
			agent.Password = "redis"
		}
		if agent.ConnectTimeout == 0 {
			agent.ConnectTimeout = 10 * time.Second
		}
		if agent.ReadTimeout == 0 {
			agent.ReadTimeout = 10 * time.Second
		}
		if agent.WriteTimeout == 0 {
			agent.WriteTimeout = 10 * time.Second
		}
		if agent.IdleTimeout == 0 {
			agent.IdleTimeout = 240 * time.Second
		}
		if agent.MaxActive == 0 {
			agent.MaxActive = 5
		}
		if agent.MaxIdle == 0 {
			agent.MaxIdle = 3
		}
		if !agent.Wait {
			agent.Wait = true
		}
		if agent.DB == 0 {
			agent.DB = 10
		}
	}
}

func (agent *Agent) GetPool() *redis.Pool {
	return agent.pool
}

func (agent *Agent) initialPool() {
	agent.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				agent.Protocol,
				fmt.Sprintf("%s:%d", agent.Host, agent.Port),
				redis.DialConnectTimeout(agent.ConnectTimeout),
				redis.DialReadTimeout(agent.ReadTimeout),
				redis.DialWriteTimeout(agent.WriteTimeout),
				redis.DialPassword(agent.Password.String()),
				redis.DialDatabase(agent.DB),
			)
		},
		MaxIdle:     agent.MaxIdle,
		MaxActive:   agent.MaxActive,
		IdleTimeout: agent.IdleTimeout,
		Wait:        true,
	}
}

func (agent *Agent) RegisterReceiver(receiver func(status *mq.TaskStatus) error) {
	agent.queue.RegisterReceiver(receiver)
}

func (agent *Agent) StartReceiver() {
	agent.queue.StartReceiver(agent.NumWorkers)
}

func (agent *Agent) StartWorker() {
	agent.queue.StartWorker(agent.Name, agent.NumWorkers)
}

func (agent *Agent) Cancel(id string) error {
	return agent.queue.Cancel(id)
}

func (agent *Agent) ListChannel() ([]string, error) {
	return agent.queue.ListChannel()
}

func (agent *Agent) ListSubject(channel string) ([]string, error) {
	return agent.queue.ListSubject(channel)
}

func SubjectAndDataFromValue(v interface{}) (string, []byte, error) {
	subject := reflectx.IndirectType(reflect.TypeOf(v)).Name()
	data, err := MarshalData(v)
	return subject, data, err
}

func (agent *Agent) SendTask(task *mq.Task) error {
	return agent.queue.SendTask(task)
}

type CronTableInfo struct {
	CronTableID string                 `json:"cronTableID"`
	Channel     string                 `json:"channel"`
	Subject     string                 `json:"subject"`
	Spec        string                 `json:"spec"`
	Args        string                 `json:"args"`
	NextTime    timelib.MySQLTimestamp `json:"nextTime"`
	Desc        string                 `json:"desc"`
}

func (agent *Agent) RegisterCron(subject string, spec string) error {
	cronTable := &CronTableInfo{}
	cronTable.CronTableID = fmt.Sprintf("%s-%s", agent.Name, subject)
	cronTable.Subject = subject
	cronTable.Channel = agent.Name
	cronTable.Spec = spec
	cronTable.Desc = cronTable.CronTableID

	bytes, _ := json.Marshal(cronTable)
	_, err := agent.Publish(agent.CentreChannel, "SubjectCreateOrUpdateCron", bytes)
	return err
}

func (agent *Agent) Next(subject string, args interface{}) (*mq.Task, error) {
	var data []byte
	if args != nil {
		data, _ = json.Marshal(args)
	}
	return agent.queue.Next(agent.Name, subject, data)
}

func (agent *Agent) Defer(subject string, nextTime time.Time, args interface{}) (*mq.Task, error) {
	cronTable := &CronTableInfo{}
	cronTable.CronTableID = uuid.New().String()
	cronTable.Channel = agent.Name
	cronTable.Subject = subject
	cronTable.NextTime = timelib.MySQLTimestamp(nextTime)

	if args != nil {
		data, _ := json.Marshal(args)
		cronTable.Args = string(data)
	}

	bytes, _ := json.Marshal(cronTable)
	return agent.Publish(agent.CentreChannel, "SubjectCreateOrUpdateCron", bytes)
}

func (agent *Agent) Publish(chanel, subject string, data []byte) (*mq.Task, error) {
	return agent.queue.Publish(chanel, subject, data)
}

func (agent *Agent) Register(subject string, job mq.Job) {
	agent.queue.Register(subject, job)
}

func (agent *Agent) RegisterRoutes(routes ...*courier.Route) {
	for _, route := range routes {
		if len(route.Operators) == 0 {
			continue
		}

		operatorMetas := courier.ToOperatorMetaList(route.Operators...)

		lastOpIndex := len(operatorMetas) - 1
		lastOpMeta := operatorMetas[lastOpIndex]
		subject := lastOpMeta.Type.Name()

		if cronDescriber, ok := lastOpMeta.Operator.(mq.CronDescriber); ok {
			spec := cronDescriber.CronSpec()
			_, err := cron.Parse(spec)
			if err != nil {
				panic(err)
			}
			if agent.RegisterCron(subject, spec); err != nil {
				panic(err)
			}
		}

		agent.queue.Register(subject, func(task *mq.Task) (interface{}, error) {
			ctx := context.Background()

			for i, opMeta := range operatorMetas {
				op := reflect.New(opMeta.Type).Interface().(courier.IOperator)

				if err := UnmarshalData(task.Data, op); err != nil {
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

func (agent *Agent) initialQueue() {
	agent.queue = mq.NewJobQueue(
		mq_redis.NewRedisBroker(agent.pool, agent.Prefix),
		mq_redis.NewRedisBackend(agent.pool, agent.Prefix),
	)
}

func (agent *Agent) Init() {
	agent.initialPool()
	agent.initialQueue()
}
