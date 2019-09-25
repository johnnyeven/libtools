package kafka

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/sirupsen/logrus"
	"github.com/wvanbergen/kazoo-go"
	"gopkg.in/Shopify/sarama.v1"

	"github.com/johnnyeven/libtools/kafka/consumergroup"
)

func PartitionerConstructor(topic string) sarama.Partitioner {
	return sarama.NewHashPartitioner(topic)
}

func (kafkaConfig KafkaConfig) initKafkaConfig(clientConfig *sarama.Config) {
	// init net config
	clientConfig.Net.DialTimeout = time.Duration(kafkaConfig.Net.DialTimeout) * time.Millisecond
	clientConfig.Net.KeepAlive = time.Duration(kafkaConfig.Net.KeepAlive) * time.Millisecond
	clientConfig.Net.MaxOpenRequests = kafkaConfig.Net.MaxOpenRequests
	clientConfig.Net.ReadTimeout = time.Duration(kafkaConfig.Net.ReadTimeout) * time.Millisecond
	clientConfig.Net.SASL.Enable = kafkaConfig.Net.SASL.Enable
	clientConfig.Net.SASL.Handshake = kafkaConfig.Net.SASL.Handshake
	clientConfig.Net.SASL.Password = kafkaConfig.Net.SASL.Password
	clientConfig.Net.SASL.User = kafkaConfig.Net.SASL.User
	clientConfig.Net.TLS.Enable = false
	clientConfig.Net.WriteTimeout = time.Duration(kafkaConfig.Net.WriteTimeout) * time.Millisecond

	// init meta config
	clientConfig.Metadata.RefreshFrequency = time.Duration(kafkaConfig.MetaData.RefreshFrequency) * time.Millisecond
	clientConfig.Metadata.Retry.Backoff = time.Duration(kafkaConfig.MetaData.Retry.Backoff) * time.Millisecond
	clientConfig.Metadata.Retry.Max = kafkaConfig.MetaData.Retry.Max

	// init producer config
	clientConfig.Producer.MaxMessageBytes = kafkaConfig.Producer.MaxMessageBytes
	clientConfig.Producer.Flush.MaxMessages = kafkaConfig.Producer.Flush.MaxMessages
	clientConfig.Producer.RequiredAcks = sarama.RequiredAcks(kafkaConfig.Producer.RequiredAcks)
	clientConfig.Producer.Retry.Backoff = time.Duration(kafkaConfig.Producer.Retry.Backoff) * time.Millisecond
	clientConfig.Producer.Retry.Max = kafkaConfig.Producer.Retry.Max
	clientConfig.Producer.Return.Errors = kafkaConfig.Producer.Return.Errors
	clientConfig.Producer.Return.Successes = kafkaConfig.Producer.Return.Successes
	clientConfig.Producer.Timeout = time.Duration(kafkaConfig.Producer.Timeout) * time.Millisecond
	clientConfig.Producer.Compression = sarama.CompressionCodec(kafkaConfig.Producer.Compression)
	clientConfig.Producer.Partitioner = PartitionerConstructor

	// init consumer config
	clientConfig.Consumer.Fetch.Default = kafkaConfig.Consumer.Fetch.Default
	clientConfig.Consumer.Fetch.Max = kafkaConfig.Consumer.Fetch.Max
	clientConfig.Consumer.Fetch.Min = kafkaConfig.Consumer.Fetch.Min
	clientConfig.Consumer.MaxProcessingTime = time.Duration(kafkaConfig.Consumer.MaxProcessingTime) * time.Millisecond
	clientConfig.Consumer.MaxWaitTime = time.Duration(kafkaConfig.Consumer.MaxWaitTime) * time.Millisecond
	clientConfig.Consumer.Offsets.CommitInterval = time.Duration(kafkaConfig.Consumer.Offsets.CommitInterval) * time.Millisecond
	clientConfig.Consumer.Offsets.Initial = kafkaConfig.Consumer.Offsets.Initial
	clientConfig.Consumer.Offsets.Retention = time.Duration(kafkaConfig.Consumer.Offsets.Retention) * time.Millisecond
	clientConfig.Consumer.Retry.Backoff = time.Duration(kafkaConfig.Consumer.Retry.Backoff) * time.Millisecond
	clientConfig.Consumer.Return.Errors = kafkaConfig.Consumer.Return.Errors

	// misc
	clientConfig.ClientID = strconv.FormatInt(int64(os.Getpid()), 10)
	clientConfig.MetricRegistry = metrics.NewRegistry()
	switch kafkaConfig.Version {
	case "0.8.2.0":
		clientConfig.Version = sarama.V0_8_2_0
	case "0.8.2.1":
		clientConfig.Version = sarama.V0_8_2_1
	case "0.8.2.2":
		clientConfig.Version = sarama.V0_8_2_2
	case "0.9.0.0":
		clientConfig.Version = sarama.V0_9_0_0
	case "0.9.0.1":
		clientConfig.Version = sarama.V0_9_0_1
	case "0.10.0.0":
		clientConfig.Version = sarama.V0_10_0_0
	case "0.10.0.1":
		clientConfig.Version = sarama.V0_10_0_1
	case "0.10.1.0":
		clientConfig.Version = sarama.V0_10_1_0
	default:
		panic(fmt.Sprintf("invalid kafka version %s", kafkaConfig.Version))
	}
}

func (kafkaConfig KafkaConfig) NewKafkaClient() (kafkaClient sarama.Client, err error) {
	clientConfig := sarama.Config{}
	kafkaConfig.initKafkaConfig(&clientConfig)
	metrics.UseNilMetrics = true

	kafkaClient, err = sarama.NewClient(kafkaConfig.Addrs, &clientConfig)
	return
}

type logger struct {
}

func (l *logger) Printf(str string, values ...interface{}) {
	logrus.WithField("tag", "kafka").Warningf("%s[%+v]", str, values)
}
func initKazooConfig(kazooConfig *kazoo.Config, zooConfig ZookeeperConfig) {
	kazooConfig.Chroot = zooConfig.Chroot
	kazooConfig.Timeout = time.Duration(zooConfig.Timeout) * time.Millisecond
	kazooConfig.Logger = &logger{}
}

func (consumerKafkaConfig ConsumerKafkaConfig) NewKafkaConsumer() (cg *consumergroup.ConsumerGroup, err error) {
	// init kafka config
	clientConfig := sarama.Config{}
	consumerKafkaConfig.ConsumerKafka.initKafkaConfig(&clientConfig)
	metrics.UseNilMetrics = true

	// init zookeeper config
	kazooConfig := kazoo.Config{}
	initKazooConfig(&kazooConfig, consumerKafkaConfig.Zookeeper)

	// init consume config
	consumeConfig := consumergroup.NewConfig()
	consumeConfig.Config = &clientConfig
	consumeConfig.Zookeeper = &kazooConfig
	consumeConfig.Offsets.CommitInterval = clientConfig.Consumer.Offsets.CommitInterval
	consumeConfig.Offsets.Initial = clientConfig.Consumer.Offsets.Initial
	consumeConfig.Offsets.ProcessingTimeout = clientConfig.Consumer.MaxProcessingTime

	cg, err = consumergroup.JoinConsumerGroup(consumerKafkaConfig.ConsumerGroupName, consumerKafkaConfig.Topics, consumerKafkaConfig.Zookeeper.Addrs, consumeConfig)
	if err != nil {
		return
	}
	return
}
func (producerKafkaConfig *ProducerKafkaConfig) Init() {
	if producerKafkaConfig.kafkaClient == nil {
		var err error
		if producerKafkaConfig.kafkaClient, err = producerKafkaConfig.ProducerKafka.NewKafkaClient(); err != nil {
			panic(fmt.Sprintf("producerKafka init fail[err:%v]", err))
		}
	}
}

func (consumerKafkaConfig *ConsumerKafkaConfig) Init() {
	if consumerKafkaConfig.kafkaClient == nil {
		var err error
		if consumerKafkaConfig.kafkaClient, err = consumerKafkaConfig.ConsumerKafka.NewKafkaClient(); err != nil {
			panic(fmt.Sprintf("Kafka client fail[err:%v]", err))
		}
	}
	if consumerKafkaConfig.cG == nil {
		var err error
		consumerKafkaConfig.cG, err = consumerKafkaConfig.NewKafkaConsumer()
		if err != nil {
			panic(fmt.Sprintf("NewKafkaConsumer fail[err:%v]", err))
		}
	}
}

func (producerKafkaConfig *ProducerKafkaConfig) GetKafkaClient() sarama.Client {
	return producerKafkaConfig.kafkaClient
}

func (consumerKafkaConfig *ConsumerKafkaConfig) GetKafkaClient() sarama.Client {
	return consumerKafkaConfig.kafkaClient
}

func (consumerKafkaConfig *ConsumerKafkaConfig) GetCG() *consumergroup.ConsumerGroup {
	return consumerKafkaConfig.cG
}
