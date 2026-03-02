package in

import (
	"context"

	"github.com/IBM/sarama"
)

type Consumer interface {
	Consume(ctx context.Context) error
}

type consumer struct {
	consumerGroup sarama.ConsumerGroup
	h             Handler
	topics        []string
}

func NewConsumer(client sarama.Client, groupID string, handler Handler) (Consumer, error) {
	cons, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		return nil, err
	}
	return &consumer{consumerGroup: cons, h: handler, topics: []string{"mytopic"}}, err
}

func (c *consumer) Consume(ctx context.Context) error {
	for {
		if err := c.consumerGroup.Consume(ctx, c.topics, c.h); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *consumer) validateEvent() error {
	return nil
}
