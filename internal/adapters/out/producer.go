package out

import (
	event "converter/internal/events"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Producer interface {
	PublishTaskCreated(event event.TaskCreated) error
}

type producer struct {
	prod  sarama.SyncProducer
	topic string
}

func NewProducer(client sarama.Client, t string) (Producer, error) {

	prod, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	return &producer{prod: prod,
		topic: t}, err
}

func (p *producer) PublishTaskCreated(event event.TaskCreated) error {
	msg := sarama.ProducerMessage{
		Topic:     p.topic,
		Key:       nil,
		Value:     nil,
		Headers:   nil,
		Metadata:  nil,
		Offset:    0,
		Partition: 0,
		Timestamp: time.Time{},
	}

	_, _, err := p.prod.SendMessage(&msg)
	if err != nil {
		return err
	}

	for {
		log.Printf("sending message %s... in topic %s:", event, p.topic)
		time.Sleep(1 * time.Second)
	}

	return nil
}
