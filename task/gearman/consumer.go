package gearman

import (
	"fmt"
	"github.com/johnnyeven/libtools/task/constants"
	"github.com/mikespook/gearman-go/worker"
	"github.com/sirupsen/logrus"
)

type GearmanConsumer struct {
	worker          *worker.Worker
	workerProcessor constants.TaskProcessor
}

func NewGearmanConsumer(info constants.ConnectionInfo) *GearmanConsumer {
	w := worker.New(worker.Unlimited)
	w.AddServer(info.Protocol, fmt.Sprintf("%s:%d", info.Host, info.Port))
	w.ErrorHandler = func(e error) {
		logrus.Errorf("worker handled err: %v", e)
	}
	return &GearmanConsumer{
		worker: w,
	}
}

func (b *GearmanConsumer) RegisterChannel(channel string, processor constants.TaskProcessor) error {
	b.workerProcessor = processor
	return b.worker.AddFunc(channel, b.processorJob, worker.Unlimited)
}

func (b *GearmanConsumer) processorJob(job worker.Job) ([]byte, error) {
	t := &constants.Task{}
	err := constants.UnmarshalData(job.Data(), t)
	if err != nil {
		return nil, err
	}
	ret, err := b.workerProcessor(t)
	if err != nil {
		return nil, err
	}
	return constants.MarshalData(ret)
}

func (b *GearmanConsumer) Work() {
	if err := b.worker.Ready(); err != nil {
		logrus.Panic("gearman worker not ready...")
	}
	logrus.Debug("GearmanConsumer.Working...")
	b.worker.Work()
}

func (b *GearmanConsumer) Stop() {
	b.worker.Close()
	logrus.Debug("GearmanConsumer.Stopped")
}
