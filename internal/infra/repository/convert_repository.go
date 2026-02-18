package repository

import (
	"converter/internal/config"
	"converter/internal/domain/entity"
	"converter/internal/domain/repository"
	"converter/internal/usecase/convert"
	"log"
	"path/filepath"
)

// ConvertRepositoryImpl - реализация ConvertRepository
type ConvertRepositoryImpl struct {
	mdFinder     convert.MDFinder
	imageFinder  convert.ImageFinder
	tmpDir       string
	converterCfg *config.ConverterConfig
}

func NewConvertRepository(
	mdFinder convert.MDFinder,
	imageFinder convert.ImageFinder,
	cfg *config.ConverterConfig,
) repository.ConvertRepository {
	return &ConvertRepositoryImpl{
		mdFinder:     mdFinder,
		imageFinder:  imageFinder,
		tmpDir:       cfg.TmpDir,
		converterCfg: cfg,
	}
}

// ScanMarkdownFiles - сканирует директорию и находит все MD файлы
func (r *ConvertRepositoryImpl) ScanMarkdownFiles(root string) (map[string]string, error) {
	return r.mdFinder.ScanFiles(root)
}

// ConvertFiles - конвертирует markdown файлы в PDF
func (r *ConvertRepositoryImpl) ConvertFiles(files map[string]string, outputDir string) (*entity.ConversionResult, error) {
	result := &entity.ConversionResult{
		Converted: make([]string, 0),
		Failed:    make([]string, 0),
	}

	if err := r.imageFinder.RecursiveFindImage(outputDir); err != nil {
		log.Printf("Warning: Failed to find images: %v", err)
	}

	processedFiles := make([]string, 0)
	for _, filePath := range files {
		tmpFile, err := r.mdFinder.ProcessFile(filePath, r.tmpDir)
		if err != nil {
			log.Printf("Failed to process file %s: %v", filePath, err)
			result.Failed = append(result.Failed, filePath)
			continue
		}
		processedFiles = append(processedFiles, tmpFile)
	}

	if len(processedFiles) == 0 {
		return result, nil
	}

	for _, file := range processedFiles {
		htmlName := changeExtensionToHTML(file)
		err := r.mdFinder.ConvertMDToPDF(file, htmlName, outputDir)
		if err != nil {
			log.Printf("Failed to convert file %s: %v", file, err)
			result.Failed = append(result.Failed, file)
			continue
		}

		result.Converted = append(result.Converted, htmlName)
		result.ProcessedCount++
	}

	return result, nil
}

// changeExtensionToHTML - меняет расширение файла с .md_TEMP.md на .html
func changeExtensionToHTML(filePath string) string {
	dir := filepath.Dir(filePath)
	file := filepath.Base(filePath)

	ext := filepath.Ext(file)
	nameWithoutExt := file[:len(file)-len(ext)]

	if len(nameWithoutExt) > 5 && nameWithoutExt[len(nameWithoutExt)-5:] == "_TEMP" {
		nameWithoutExt = nameWithoutExt[:len(nameWithoutExt)-5]
	}

	return filepath.Join(dir, nameWithoutExt+".html")
}
