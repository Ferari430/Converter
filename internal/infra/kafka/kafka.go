package kafka

import (
	"converter/internal/config"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

type KafkaClient struct {
	client sarama.Client
	config config.KafkaConfig
}

func NewClient(kafkaCfg config.KafkaConfig) (*KafkaClient, error) {
	log.Printf("Connecting to Kafka at: %s", kafkaCfg.BrokersAddr)

	if kafkaCfg.BrokersAddr == "" {
		return nil, fmt.Errorf("Kafka broker address is empty")
	}

	if !strings.Contains(kafkaCfg.BrokersAddr, ":") {
		return nil, fmt.Errorf("invalid Kafka address format: %s, expected host:port", kafkaCfg.BrokersAddr)
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V3_6_0_0
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaCfg.Producer.Return.Successes = true

	saramaCfg.Net.DialTimeout = 5 * time.Second
	saramaCfg.Net.ReadTimeout = 5 * time.Second
	saramaCfg.Net.WriteTimeout = 5 * time.Second
	saramaCfg.Net.KeepAlive = 5 * time.Second

	addr := []string{kafkaCfg.BrokersAddr}

	log.Printf("Trying to connect to Kafka brokers: %v", addr)

	client, err := sarama.NewClient(addr, saramaCfg)
	if err != nil {
		log.Printf("Failed to connect to Kafka: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to Kafka")

	brokers := client.Brokers()
	log.Printf("Available brokers: %d", len(brokers))
	for _, broker := range brokers {
		log.Printf("Broker: %s", broker.Addr())
	}

	topics, err := client.Topics()
	if err != nil {
		log.Printf("Warning: cannot list topics: %v", err)
		return nil, err
	}

	for _, topic := range topics {
		log.Printf("Topic: %s", topic)
	}

	return &KafkaClient{
		client: client,
		config: kafkaCfg,
	}, nil
}

func (kc *KafkaClient) Producer() (sarama.SyncProducer, error) {
	return sarama.NewSyncProducerFromClient(kc.client)
}

func (kc *KafkaClient) ConsumerGroup() (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroupFromClient(kc.config.ConsumerGroupID, kc.client)
}

func (kc *KafkaClient) Close() error {
	return kc.client.Close()
}

func NewConsumer(
	client *KafkaClient,
	groupID string,
	handler sarama.ConsumerGroupHandler,
	topic string,
) (Consumer, error) {
	consumerGroup, err := client.ConsumerGroup()
	if err != nil {
		return nil, err
	}

	return NewConsumerImpl(consumerGroup, handler, topic), nil
}

func stringPtr(s string) *string {
	return &s
}
