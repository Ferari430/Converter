package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type Consumer interface {
	Consume(ctx context.Context) error
	Close() error
}

type ConsumerImpl struct {
	consumerGroup sarama.ConsumerGroup
	handler       sarama.ConsumerGroupHandler
	topic         string
}

func NewConsumerImpl(
	consumerGroup sarama.ConsumerGroup,
	handler sarama.ConsumerGroupHandler,
	topic string,
) Consumer {
	return &ConsumerImpl{
		consumerGroup: consumerGroup,
		handler:       handler,
		topic:         topic,
	}
}

func (c *ConsumerImpl) Consume(ctx context.Context) error {
	return c.consumerGroup.Consume(ctx, []string{c.topic}, c.handler)
}

func (c *ConsumerImpl) Close() error {
	return c.consumerGroup.Close()
}
