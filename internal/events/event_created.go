package event

import "time"

type TaskCreated struct {
	EventID   string    `json:"event_id"`   // Уникальный ID события
	EventType string    `json:"event_type"` // "task_created"
	TaskID    string    `json:"task_id"`    // Уникальный ID задачи
	ChatID    int64     `json:"chat_id"`    // Telegram Chat ID
	FilePath  string    `json:"file_path"`  // Путь к архиву
	FileName  string    `json:"file_name"`  // Имя файла
	CreatedAt time.Time `json:"created_at"` // Время создания
}
