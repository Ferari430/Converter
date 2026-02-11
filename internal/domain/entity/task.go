package entity

import "time"

// ConversionTask - доменная сущность для задачи конвертирования
type ConversionTask struct {
	TaskID    string
	ChatID    int64
	FilePath  string // путь к ZIP архиву
	CreatedAt time.Time
}

// ConversionResult - результат конвертирования
type ConversionResult struct {
	TaskID         string
	SourceDir      string   // директория с распакованными файлами
	Converted      []string // пути к конвертированным файлам
	Failed         []string // файлы которые не удалось конвертировать
	ProcessedCount int64
}
