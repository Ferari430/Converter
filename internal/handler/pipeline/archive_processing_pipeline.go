package pipeline

import (
	"converter/internal/handler/convert"
	"converter/internal/handler/unzip"
	"log"
)

// ArchiveProcessingPipeline объединяет все стадии обработки архива
type ArchiveProcessingPipeline struct {
	unzipHandler   *unzip.UnzipHandler
	convertHandler *convert.ConvertHandler
	tmpDir         string
}

func NewArchiveProcessingPipeline(
	uh *unzip.UnzipHandler,
	ch *convert.ConvertHandler,
	tmpDir string,
) *ArchiveProcessingPipeline {
	return &ArchiveProcessingPipeline{
		unzipHandler:   uh,
		convertHandler: ch,
		tmpDir:         tmpDir,
	}
}

// Process запускает полный pipeline обработки архива
// Стадия 1: Распаковывает архив
// Стадия 2: Обрабатывает распакованные файлы
func (app *ArchiveProcessingPipeline) Process(archivePath string) error {
	log.Println("[Pipeline] Starting archive processing for:", archivePath)

	// Стадия 1: Распаковка архива
	log.Println("[Pipeline] Stage 1: Unzipping archive...")
	unzippedDir, err := app.unzipHandler.Unzip(archivePath)
	if err != nil {
		log.Println("[Pipeline] Error unzipping archive:", err)
		return err
	}

	log.Println("[Pipeline] Archive unzipped to:", unzippedDir)

	// Стадия 2: Обработка распакованных файлов
	log.Println("[Pipeline] Stage 2: Processing files...")
	app.convertHandler.HandleDirPipline(unzippedDir, app.tmpDir)
	log.Println("[Pipeline] Archive processing completed")

	return nil
}
