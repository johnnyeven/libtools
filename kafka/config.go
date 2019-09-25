package kafka

import (
	"fmt"

	"gopkg.in/Shopify/sarama.v1"

	"github.com/johnnyeven/libtools/conf"
	"github.com/johnnyeven/libtools/kafka/consumergroup"
)

type ProducerKafkaConfig struct {
	ProducerKafka KafkaConfig
	kafkaClient   sarama.Client
}

func (producerKafkaConfig ProducerKafkaConfig) MarshalDefaults(v interface{}) {
	if c, ok := v.(*ProducerKafkaConfig); ok {
		if c.ProducerKafka.Name == "" {
			c.ProducerKafka.Name = "producekafka"
		}
	}
}

type ConsumerKafkaConfig struct {
	ConsumerKafka     KafkaConfig
	Zookeeper         ZookeeperConfig
	Topics            []string
	ConsumerGroupName string
	cG                *consumergroup.ConsumerGroup
	kafkaClient       sarama.Client
}

func (consumerKafkaConfig ConsumerKafkaConfig) MarshalDefaults(v interface{}) {
	if c, ok := v.(*ConsumerKafkaConfig); ok {
		if c.ConsumerKafka.Name == "" {
			c.ConsumerKafka.Name = "consumekafka"
		}
		if c.Zookeeper.Name == "" {
			c.Zookeeper.Name = "consumekafka"
		}
	}
}

type KafkaConfig struct {
	Name     string
	Addrs    []string `conf:"env"`
	Net      NetConfig
	MetaData MetaDataConfig
	Producer ProducerConfig
	Consumer ConsumerConfig
	Version  string
}

func (kafkaConfig KafkaConfig) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Addrs": []string{
			fmt.Sprintf("%s:%d",
				conf.RancherInternal("tool-deps", kafkaConfig.Name).String(),
				9092,
			),
		},
	}
}

func (kafkaConfig KafkaConfig) MarshalDefaults(v interface{}) {
	if c, ok := v.(*KafkaConfig); ok {
		if c.Version == "" {
			c.Version = "0.8.2.2"
		}
	}
}

type ZookeeperConfig struct {
	Name    string
	Addrs   []string `conf:"env"`
	Timeout int64
	Chroot  string
}

func (zookeeperConfig ZookeeperConfig) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Addrs": []string{
			fmt.Sprintf("%s:%d",
				conf.RancherInternal("tool-deps", zookeeperConfig.Name).String(),
				2181,
			),
		},
	}
}

func (zookeeperConfig ZookeeperConfig) MarshalDefaults(v interface{}) {
	if c, ok := v.(*ZookeeperConfig); ok {
		if c.Timeout == 0 {
			c.Timeout = 2000
		}
	}
}

type NetConfig struct {
	MaxOpenRequests int
	DialTimeout     int64
	ReadTimeout     int64
	WriteTimeout    int64
	KeepAlive       int64
	SASL            SASLConfig
}

func (netConfig NetConfig) MarshalDefaults(v interface{}) {
	if c, ok := v.(*NetConfig); ok {
		if c.MaxOpenRequests == 0 {
			c.MaxOpenRequests = 3
		}
		if c.DialTimeout == 0 {
			c.DialTimeout = 500
		}
		if c.ReadTimeout == 0 {
			c.ReadTimeout = 2000
		}
		if c.WriteTimeout == 0 {
			c.WriteTimeout = 1000
		}
	}
}

type SASLConfig struct {
	Enable    bool
	Handshake bool
	User      string
	Password  string
}

type RetryConfig struct {
	Max     int
	Backoff int64
}

func (retryConfig RetryConfig) MarshalDefaults(v interface{}) {
	if c, ok := v.(*RetryConfig); ok {
		if c.Backoff == 0 {
			c.Backoff = 200
		}
	}
}

type MetaDataConfig struct {
	Retry            RetryConfig
	RefreshFrequency int64
}

func (metaDataConfig MetaDataConfig) MarshalDefaults(v interface{}) {
	if c, ok := v.(*MetaDataConfig); ok {
		if c.RefreshFrequency == 0 {
			c.RefreshFrequency = 2000
		}
	}
}

type ReturnConfig struct {
	Successes bool
	Errors    bool
}

func (returnConfig ReturnConfig) MarshalDefaults(v interface{}) {
	if c, ok := v.(*ReturnConfig); ok {
		if !c.Successes {
			c.Successes = true
		}
		if !c.Errors {
			c.Errors = true
		}
	}
}

type FlushConfig struct {
	MaxMessages int
}

func (flushConfig FlushConfig) MarshalDefaults(v interface{}) {
	if c, ok := v.(*FlushConfig); ok {
		if c.MaxMessages == 0 {
			c.MaxMessages = 1
		}
	}
}

type ProducerConfig struct {
	MaxMessageBytes int
	RequiredAcks    int
	Timeout         int64
	Compression     int8
	Return          ReturnConfig
	Flush           FlushConfig
	Retry           RetryConfig
}

func (producerConfig ProducerConfig) MarshalDefaults(v interface{}) {
	if c, ok := v.(*ProducerConfig); ok {
		if c.MaxMessageBytes == 0 {
			c.MaxMessageBytes = 1024000
		}
		if c.RequiredAcks == 0 {
			c.RequiredAcks = 1
		}
		if c.Timeout == 0 {
			c.Timeout = 2000
		}
	}
}

type FetchConfig struct {
	Min     int32
	Default int32
	Max     int32
}

func (fetchConfig FetchConfig) MarshalDefaults(v interface{}) {
	if c, ok := v.(*FetchConfig); ok {
		if c.Min == 0 {
			c.Min = 1
		}
		if c.Default == 0 {
			c.Default = 32678
		}
	}
}

type OffsetsConfig struct {
	CommitInterval int64
	Initial        int64
	Retention      int64
}

func (offsetsConfig OffsetsConfig) MarshalDefaults(v interface{}) {
	if c, ok := v.(*OffsetsConfig); ok {
		if c.CommitInterval == 0 {
			c.CommitInterval = 1000
		}
		if c.Initial == 0 {
			c.Initial = -2
		}
		if c.Retention == 0 {
			c.Retention = 3600000
		}
	}
}

type ConsumerConfig struct {
	Retry             RetryConfig
	Fetch             FetchConfig
	MaxWaitTime       int64
	MaxProcessingTime int64
	Return            ReturnConfig
	Offsets           OffsetsConfig
}

func (consumerConfig ConsumerConfig) MarshalDefaults(v interface{}) {
	if c, ok := v.(*ConsumerConfig); ok {
		if c.MaxWaitTime == 0 {
			c.MaxWaitTime = 200
		}
		if c.MaxProcessingTime == 0 {
			c.MaxProcessingTime = 1000
		}
	}
}
