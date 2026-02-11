package repository

import "converter/internal/domain/entity"

// UnzipRepository - интерфейс для работы с распаковкой архивов
type UnzipRepository interface {
	Unzip(archivePath string) (*entity.UnzipResult, error)
}

// ConvertRepository - интерфейс для работы с конвертированием
type ConvertRepository interface {
	ScanMarkdownFiles(root string) (map[string]string, error)
	ConvertFiles(files map[string]string, outputDir string) (*entity.ConversionResult, error)
}

// EventPublisher - интерфейс для публикации событий в очередь
type EventPublisher interface {
	PublishConverted(result *entity.ConversionResult) error
}
