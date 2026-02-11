package usecase

import (
	"context"
	"fmt"
	"log"

	"converter/internal/domain/entity"
	"converter/internal/domain/repository"
)

// ConvertArchiveUseCase - основной Use Case для обработки архива
type ConvertArchiveUseCase struct {
	unzipRepo      repository.UnzipRepository
	convertRepo    repository.ConvertRepository
	eventPublisher repository.EventPublisher
}

func NewConvertArchiveUseCase(
	unzipRepo repository.UnzipRepository,
	convertRepo repository.ConvertRepository,
	eventPublisher repository.EventPublisher,
) *ConvertArchiveUseCase {
	return &ConvertArchiveUseCase{
		unzipRepo:      unzipRepo,
		convertRepo:    convertRepo,
		eventPublisher: eventPublisher,
	}
}

// Execute - основной метод для выполнения задачи конвертирования
func (uc *ConvertArchiveUseCase) Execute(ctx context.Context, task *entity.ConversionTask) error {
	log.Printf("Starting conversion task: TaskID=%s, FilePath=%s", task.TaskID, task.FilePath)

	unzipResult, err := uc.unzipRepo.Unzip(task.FilePath)
	if err != nil {
		log.Printf("Failed to unzip archive: %v", err)
		return fmt.Errorf("unzip failed: %w", err)
	}

	log.Printf("Archive unpacked successfully. Files: %d, DestDir: %s", unzipResult.Count, unzipResult.DestDir)

	mdFiles, err := uc.convertRepo.ScanMarkdownFiles(unzipResult.DestDir)
	if err != nil {
		log.Printf("Failed to scan markdown files: %v", err)
		return fmt.Errorf("scan files failed: %w", err)
	}

	if len(mdFiles) == 0 {
		log.Printf("No markdown files found in: %s", unzipResult.DestDir)
		return fmt.Errorf("no markdown files found")
	}

	log.Printf("Found %d markdown files", len(mdFiles))

	result, err := uc.convertRepo.ConvertFiles(mdFiles, unzipResult.DestDir)
	if err != nil {
		log.Printf("Failed to convert files: %v", err)
		return fmt.Errorf("conversion failed: %w", err)
	}

	result.TaskID = task.TaskID
	result.SourceDir = unzipResult.DestDir

	log.Printf("Conversion completed. Processed: %d, Failed: %d", result.ProcessedCount, len(result.Failed))

	if err := uc.eventPublisher.PublishConverted(result); err != nil {
		log.Printf("Failed to publish event: %v", err)
	}

	return nil
}
