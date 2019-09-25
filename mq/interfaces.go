package mq

type Broker interface {
	SendTask(task *Task) (err error)
	GetTask(channel string) (task *Task, err error)
}

type Canceller interface {
	Cancel(id string) (err error)
	IsCancelled(id string) (ok bool, err error)
	ClearCancellation(id string) (err error)
}

type ChannelMgr interface {
	ListChannel() (channelList []string, err error)
	ListSubject(channel string) (subjectList []string, err error)
	RegisterChannel(channel string, subjectList []string) error
}

type Backend interface {
	Canceller
	ChannelMgr
	FeedBack(taskStatus *TaskStatus) (err error)
	GetFeedback() (taskStatus *TaskStatus, err error)
}

type CronDescriber interface {
	CronSpec() string
}

type MsgHead struct {
	Channel string `json:"channel"`
	Subject string `json:"subject"`
	ID      string `json:"id"`
}
