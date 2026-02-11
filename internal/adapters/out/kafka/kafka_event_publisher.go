package kafka

import (
	"converter/internal/domain/entity"
	"converter/internal/domain/repository"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

// ConvertedEventDTO - DTO для отправки события о завершении конвертирования
type ConvertedEventDTO struct {
	EventID     string   `json:"event_id"`
	EventType   string   `json:"event_type"`
	TaskID      string   `json:"task_id"`
	SourceDir   string   `json:"source_dir"`
	Converted   []string `json:"converted"`
	Failed      []string `json:"failed"`
	Processed   int64    `json:"processed"`
	CompletedAt string   `json:"completed_at"`
}

// KafkaEventPublisher - публикует события в Kafka
type KafkaEventPublisher struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaEventPublisher(producer sarama.SyncProducer, topic string) repository.EventPublisher {
	return &KafkaEventPublisher{
		producer: producer,
		topic:    topic,
	}
}

// PublishConverted - публикует событие о завершении конвертирования
func (p *KafkaEventPublisher) PublishConverted(result *entity.ConversionResult) error {
	event := ConvertedEventDTO{
		EventID:   "converted_" + result.TaskID,
		EventType: "task_converted",
		TaskID:    result.TaskID,
		SourceDir: result.SourceDir,
		Converted: result.Converted,
		Failed:    result.Failed,
		Processed: result.ProcessedCount,
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(result.TaskID),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	log.Printf("Event published: Topic=%s, Partition=%d, Offset=%d", p.topic, partition, offset)
	return nil
}
