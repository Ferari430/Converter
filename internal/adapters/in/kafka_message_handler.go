package in

import (
	"context"
	"converter/internal/domain/entity"
	"converter/internal/usecase"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

// ConversionEventDTO - DTO для сообщения из Kafka
type ConversionEventDTO struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	TaskID    string `json:"task_id"`
	ChatID    int64  `json:"chat_id"`
	FilePath  string `json:"file_path"`
	CreatedAt string `json:"created_at"`
}

// KafkaMessageHandler - обработчик сообщений Kafka
type KafkaMessageHandler struct {
	useCase *usecase.ConvertArchiveUseCase
}

func NewKafkaMessageHandler(useCase *usecase.ConvertArchiveUseCase) *KafkaMessageHandler {
	return &KafkaMessageHandler{
		useCase: useCase,
	}
}

// Setup - реализация интерфейса sarama.ConsumerGroupHandler
func (h *KafkaMessageHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup - реализация интерфейса sarama.ConsumerGroupHandler
func (h *KafkaMessageHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim - реализация интерфейса sarama.ConsumerGroupHandler
func (h *KafkaMessageHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for message := range claim.Messages() {
		if err := h.handleMessage(message); err != nil {
			log.Printf("Failed to handle message: %v", err)
			// Продолжаем обработку, не прерываем consumer
		}

		session.MarkMessage(message, "")
	}
	return nil
}

// handleMessage - обработка одного сообщения
func (h *KafkaMessageHandler) handleMessage(msg *sarama.ConsumerMessage) error {
	log.Printf("Received Kafka message: Topic=%s, Partition=%d, Offset=%d", msg.Topic, msg.Partition, msg.Offset)

	var event ConversionEventDTO
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return err
	}

	log.Printf("Parsed event: TaskID=%s, FilePath=%s, EventType=%s", event.TaskID, event.FilePath, event.EventType)

	if event.EventType != "task_created" {
		log.Printf("Skipping event with type '%s' (only processing 'task_created'): TaskID=%s", event.EventType, event.TaskID)
		return nil
	}

	task := &entity.ConversionTask{
		TaskID:    event.TaskID,
		ChatID:    event.ChatID,
		FilePath:  event.FilePath,
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := h.useCase.Execute(ctx, task); err != nil {
		log.Printf("Use case execution failed: %v", err)
		return err
	}

	log.Printf("Task completed successfully: TaskID=%s", event.TaskID)
	return nil
}
