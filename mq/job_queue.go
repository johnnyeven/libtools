package mq

import (
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"
)

func NewJobQueue(broker Broker, backend Backend) *JobQueue {
	return &JobQueue{
		broker:  broker,
		backend: backend,
		jobs:    sync.Map{},
	}
}

type JobQueue struct {
	broker  Broker
	backend Backend

	channel string
	jobs    sync.Map
	worker  *Worker

	receiver       func(*TaskStatus) error
	receiverWorker *Worker
}

func (mq *JobQueue) ListChannel() ([]string, error) {
	return mq.backend.ListChannel()
}

func (mq *JobQueue) ListSubject(channel string) ([]string, error) {
	return mq.backend.ListSubject(channel)
}

func (mq *JobQueue) RegisterReceiver(receiver func(*TaskStatus) error) {
	mq.receiver = receiver
}

func (mq *JobQueue) StartReceiver(numWorkers int) {
	mq.receiverWorker = NewWorker(mq.receiverProcess, numWorkers)
	mq.receiverWorker.Start()
}

func (mq *JobQueue) StopReceiver() {
	if mq.receiverWorker != nil {
		mq.receiverWorker.Stop()
	}
}

func (mq *JobQueue) receiverProcess() error {
	taskStatus, err := mq.backend.GetFeedback()
	if err != nil || taskStatus == nil {
		return nil
	}

	if taskStatus.Status == StatusCancelled {
		if err := mq.backend.ClearCancellation(taskStatus.ID); err != nil {
			logrus.Warnf("clear cancellation %s failed %s", taskStatus.ID, err)
		}
	}

	if mq.receiver != nil {
		if err := mq.receiver(taskStatus); err != nil {
			return err
		}
	}

	return nil
}

type Job = func(task *Task) (result interface{}, err error)

func (mq *JobQueue) Register(subject string, job Job) {
	mq.jobs.Store(subject, job)
}

func (mq *JobQueue) Cancel(id string) error {
	logrus.Debugf("cancelling %s", id)
	return mq.backend.Cancel(id)
}

func (mq *JobQueue) SendTask(task *Task) error {
	return mq.broker.SendTask(task)
}

func (mq *JobQueue) Next(channel string, subject string, data []byte) (*Task, error) {
	task := NewTask(channel, subject, data)
	return task, mq.backend.FeedBack(task.Pending())
}

func (mq *JobQueue) Publish(channel string, subject string, data []byte) (*Task, error) {
	task := NewTask(channel, subject, data)
	return task, mq.SendTask(task)
}

func (mq *JobQueue) StartWorker(channel string, numWorkers int) {
	subjects := make([]string, 0)
	mq.jobs.Range(func(key, value interface{}) bool {
		subjects = append(subjects, key.(string))
		return true
	})

	mq.channel = channel

	if err := mq.backend.RegisterChannel(mq.channel, subjects); err != nil {
		logrus.Panic(err)
	}

	mq.worker = NewWorker(mq.jobProcess, numWorkers)
	mq.worker.Start()
}

func (mq *JobQueue) StopWorker() {
	if mq.worker != nil {
		mq.worker.Stop()
	}
}

func (mq *JobQueue) jobProcess() error {
	task, err := mq.broker.GetTask(mq.channel)
	if err != nil || task == nil {
		return nil
	}

	defer func() {
		if e := recover(); e != nil {
			jobErr := fmt.Errorf("panic: %s; calltrace:%s", fmt.Sprint(e), string(debug.Stack()))
			logrus.Warnf("%s.%s failed: %s", task.Subject, task.ID, jobErr)
			if err := mq.backend.FeedBack(task.Failed(jobErr)); err != nil {
				logrus.Errorf("feed back FAILED failed %s/%s.%s: %s", task.Channel, task.Subject, task.ID, err)
			}
		}
	}()

	if isCancelled, err := mq.backend.IsCancelled(task.ID); isCancelled {
		if err := mq.backend.FeedBack(task.Cancelled()); err != nil {
			logrus.Errorf("feed back CANCELLED failed %s/%s.%s: %s", task.Channel, task.Subject, task.ID, err)
		}
		return err
	}

	job, ok := mq.jobs.Load(task.Subject)
	if !ok {
		jobErr := fmt.Errorf("missing subject %s", task)
		if err := mq.backend.FeedBack(task.Failed(jobErr)); err != nil {
			logrus.Errorf("feed back FAILED failed %s/%s.%s: %s", task.Channel, task.Subject, task.ID, err)
		}
		return jobErr
	}

	if err := mq.backend.FeedBack(task.Running()); err != nil {
		logrus.Errorf("feed back RUNNING failed %s/%s.%s: %s", task.Channel, task.Subject, task.ID, err)
	}

	jobResult, jobErr := job.(Job)(task)

	if jobErr != nil {
		logrus.Warnf("%s.%s failed: %s", task.Subject, task.ID, jobErr)
		if err := mq.backend.FeedBack(task.Failed(jobErr)); err != nil {
			logrus.Errorf("feed back FAILED failed %s/%s.%s: %s", task.Channel, task.Subject, task.ID, err)
		}
		return jobErr
	}

	if err := mq.backend.FeedBack(task.Success(jobResult)); err != nil {
		logrus.Errorf("feed back SUCCESS failed %s/%s.%s: %s", task.Channel, task.Subject, task.ID, err)
	}

	return nil
}
