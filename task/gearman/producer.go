package gearman

import (
	"fmt"
	"github.com/johnnyeven/libtools/task/constants"
	"github.com/mikespook/gearman-go/client"
	"github.com/sirupsen/logrus"
)

type GearmanProducer struct {
	client *client.Client
}

func NewGearmanProducer(info constants.ConnectionInfo) *GearmanProducer {
	c, err := client.New(info.Protocol, fmt.Sprintf("%s:%d", info.Host, info.Port))
	if err != nil {
		logrus.Panicf("NewGearmanProducer err: %v", err)
	}

	return &GearmanProducer{
		client: c,
	}
}

func (c *GearmanProducer) SendTask(task *constants.Task) error {
	data, err := constants.MarshalData(task)
	if err != nil {
		logrus.Errorf("GearmanProducer.SendTask err: %v", err)
		return err
	}
	_, err = c.client.DoBg(task.Channel, data, client.JobNormal)
	if err != nil {
		return err
	}

	return nil
}

func (c *GearmanProducer) Stop() {
	c.client.Close()
}
