package kafka

import (
	"context"
	"fmt"
	"github.com/johnnyeven/libtools/task/constants"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type KafkaProducer struct {
	w *kafka.Writer
}

func NewKafkaProducer(info constants.ConnectionInfo) *KafkaProducer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:           []string{fmt.Sprintf("%s:%d", info.Host, info.Port)},
		Topic:             "test",
		Balancer:          &kafka.LeastBytes{},
	})
	return &KafkaProducer{
		w: w,
	}
}

func (p *KafkaProducer) SendTask(task *constants.Task) error {
	data, err := constants.MarshalData(task)
	if err != nil {
		logrus.Errorf("GearmanProducer.SendTask err: %v", err)
		return err
	}
	err = p.w.WriteMessages(context.Background(), kafka.Message{
		Topic: task.Channel,
		Key:   []byte(task.Subject),
		Value: data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *KafkaProducer) Stop() {
	p.w.Close()
}
