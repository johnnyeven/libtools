package kafka

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"gopkg.in/Shopify/sarama.v1"
)

func SendNotify(client sarama.Client, topic, key string, notify interface{}) error {
	valueBytes, err := json.Marshal(notify)
	if err != nil {
		logrus.Errorf("marshal pay notify request to cashdesk failed![err:%s]", err.Error())
		return err
	}
	message := sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(valueBytes),
	}

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		logrus.Errorf("new producer from client failed![err:%s]", err.Error())
		return err
	}

	defer func() {
		if err := producer.Close(); err != nil {
			logrus.Errorf("failed to close producer![err:%s]", err.Error())
		}
	}()

	partition, offset, err := producer.SendMessage(&message)
	if err != nil {
		logrus.Errorf("failed to send message![err:%s]", err.Error())
		return err
	}
	logrus.Debugf("Kafak Send Topic[%s] Message:[%s]", topic, message)
	logrus.Infof("[key:%s]send message to kafka succeed![partition:%d][offset:%d][topic:%s]", key, partition, offset, topic)
	return nil
}
