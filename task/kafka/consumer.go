package kafka

import (
	"context"
	"fmt"
	"github.com/johnnyeven/libtools/task/constants"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"time"
)

type KafkaConsumer struct {
	r               *kafka.Reader
	config          kafka.ReaderConfig
	workerProcessor constants.TaskProcessor
}

func NewKafkaConsumer(info constants.ConnectionInfo) *KafkaConsumer {
	c := kafka.ReaderConfig{
		Brokers:  []string{fmt.Sprintf("%s:%d", info.Host, info.Port)},
		MinBytes: 10e3,
		MaxBytes: 10e6,
		CommitInterval: time.Second,
	}

	return &KafkaConsumer{
		config: c,
	}
}

func (c *KafkaConsumer) RegisterChannel(channel string, processor constants.TaskProcessor) error {
	c.config.Topic = channel
	c.config.GroupID = channel + "_GROUP"
	c.workerProcessor = processor

	r := kafka.NewReader(c.config)
	c.r = r

	return nil
}

func (c *KafkaConsumer) Work() {
	logrus.Debug("KafkaConsumer.Working...")
	ctx := context.Background()
	for {
		m, err := c.r.ReadMessage(ctx)
		if err != nil {
			logrus.Errorf("kafka.Reader.ReadMessage err: %v", err)
			break
		}
		c.r.CommitMessages(ctx, m)

		t := &constants.Task{}
		err = constants.UnmarshalData(m.Value, t)
		if err != nil {
			logrus.Errorf("Work UnmarshalData err: %v", err)
			continue
		}
		_, err = c.workerProcessor(t)
		if err != nil {
			continue
		}
	}
}

func (c *KafkaConsumer) Stop() {
	c.r.Close()
	logrus.Debug("KafkaConsumer.Stopped")
}
