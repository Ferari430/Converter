package kafka

import (
	event "converter/internal/events"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
)

func TaskCreatedToMessage(topic string, task event.TaskCreated) (*sarama.ProducerMessage, error) {
	bytes, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(task.TaskID), // Key = TaskID
		Value: sarama.ByteEncoder(bytes),
	}, nil
}

func MessageToTaskCreated(msg *sarama.ConsumerMessage) (*event.TaskCreated, error) {
	var task event.TaskCreated
	if err := json.Unmarshal(msg.Value, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func TaskConvertedToMessage(topic string, task event.TaskConverted) (*sarama.ProducerMessage, error) {
	bytes, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(task.TaskID), // Key = TaskID
		Value: sarama.ByteEncoder(bytes),
	}, nil
}

func MessageToTaskConverted(msg *sarama.ConsumerMessage) (*event.TaskConverted, error) {
	var task event.TaskConverted
	if err := json.Unmarshal(msg.Value, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func KafkaMessageToConsumerMessage(domainMsg *event.KafkaMessage) *sarama.ConsumerMessage {
	return &sarama.ConsumerMessage{
		Topic:     domainMsg.Topic,
		Partition: domainMsg.Partition,
		Offset:    domainMsg.Offset,
		Key:       domainMsg.Key,
		Value:     domainMsg.Value,
		Headers:   nil,
		Timestamp: time.Now(), // или берем из заголовков
	}
}

func ConsumerMessageToKafkaMessage(saramaMsg *sarama.ConsumerMessage) *event.KafkaMessage {
	return &event.KafkaMessage{
		Topic:     saramaMsg.Topic,
		Partition: saramaMsg.Partition,
		Offset:    saramaMsg.Offset,
		Key:       saramaMsg.Key,
		Value:     saramaMsg.Value,
		Headers:   nil,
	}
}
