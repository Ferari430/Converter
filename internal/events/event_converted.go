package event

import (
	"time"
)

type KafkaMessage struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       []byte
	Value     []byte
	Headers   map[string]string
}

type EventHandler interface {
	Handle(msg *KafkaMessage) error
}

// TaskConverted - твой бизнес-ивент (переносим из твоего пакета event)
type TaskConverted struct {
	EventID   string    `json:"event_id"`   // Уникальный ID события
	EventType string    `json:"event_type"` // "task_created"
	TaskID    string    `json:"task_id"`    // Уникальный ID задачи
	ChatID    int64     `json:"chat_id"`    // Telegram Chat ID
	FilePath  string    `json:"file_path"`  // Путь к архиву
	DirName   string    `json:"dir_name"`   // Имя файла
	CreatedAt time.Time `json:"created_at"` // Время создания
}
