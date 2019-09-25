package task

import "github.com/johnnyeven/libtools/task/constants"

type Consumer interface {
	RegisterChannel(channel string, processor constants.TaskProcessor) error
	Work()
	Stop()
}

type Producer interface {
	SendTask(task *constants.Task) error
	Stop()
}

type CronDescriber interface {
	CronSpec() string
}
