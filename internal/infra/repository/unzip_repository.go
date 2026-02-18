package repository

import (
	"converter/internal/domain/entity"
	"converter/internal/domain/repository"
	"converter/internal/usecase/unzip"
)

// UnzipRepositoryImpl - реализация UnzipRepository
type UnzipRepositoryImpl struct {
	unzipService *unzip.UnzipService
}

func NewUnzipRepository(unzipService *unzip.UnzipService) repository.UnzipRepository {
	return &UnzipRepositoryImpl{
		unzipService: unzipService,
	}
}

// Unzip - распакует архив и вернет результат
func (r *UnzipRepositoryImpl) Unzip(archivePath string) (*entity.UnzipResult, error) {
	result, err := r.unzipService.UnzipArchive(archivePath)
	if err != nil {
		return nil, err
	}

	return &entity.UnzipResult{
		DestDir: result.DestDir,
		Files:   result.Files,
		Count:   result.Count,
	}, nil
}
