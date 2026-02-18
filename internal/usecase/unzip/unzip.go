package unzip

import (
	"archive/zip"
	"converter/internal/domain/entity"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type UnzipService struct{}

func NewUnzipService() *UnzipService {
	return &UnzipService{}
}

func (u *UnzipService) UnzipArchive(src string) (*entity.UnzipResult, error) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", src)
	}

	destDir, err := u.createUniqueDestDir(src)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	r, err := zip.OpenReader(src)
	if err != nil {
		return &entity.UnzipResult{DestDir: destDir}, fmt.Errorf("failed to open zip file: %w", err)
	}
	defer r.Close()

	extractedFiles, err := u.extractAllFiles(r, destDir)
	if err != nil {
		return &entity.UnzipResult{DestDir: destDir, Files: extractedFiles}, err
	}

	return &entity.UnzipResult{
		DestDir: destDir,
		Files:   extractedFiles,
		Count:   len(extractedFiles),
	}, nil
}

func (u *UnzipService) createUniqueDestDir(src string) (string, error) {
	baseDir := filepath.Dir(src)
	baseName := strings.TrimSuffix(filepath.Base(src), filepath.Ext(src))

	destDir := filepath.Join(baseDir, baseName)

	counter := 1
	originalDestDir := destDir
	for {
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			break
		}
		destDir = fmt.Sprintf("%s_%d", originalDestDir, counter)
		counter++
		if counter > 100 {
			return os.MkdirTemp("", baseName+"_*")
		}
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}

	return destDir, nil
}

// extractAllFiles распаковывает все файлы из архива
func (u *UnzipService) extractAllFiles(r *zip.ReadCloser, destDir string) ([]string, error) {
	var extractedFiles []string

	for _, f := range r.File {
		filePath, err := u.extractSingleFile(f, destDir)
		if err != nil {
			return extractedFiles, fmt.Errorf("failed to extract %s: %w", f.Name, err)
		}
		extractedFiles = append(extractedFiles, filePath)
	}

	return extractedFiles, nil
}

// extractSingleFile извлекает один файл из архива
func (u *UnzipService) extractSingleFile(f *zip.File, destDir string) (string, error) {
	// Безопасно вычисляем путь
	path := filepath.Join(destDir, f.Name)

	// Защита от ZipSlip
	if !u.isSafePath(path, destDir) {
		return "", fmt.Errorf("unsafe file path: %s", f.Name)
	}

	// Создаем директорию
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(path, f.Mode()); err != nil {
			return "", err
		}
		return path, nil
	}

	// Создаем родительские директории
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}

	// Открываем и копируем файл
	if err := u.copyFile(f, path); err != nil {
		return "", err
	}

	return path, nil
}

// isSafePath проверяет безопасность пути (защита от ZipSlip)
func (u *UnzipService) isSafePath(path, destDir string) bool {
	cleanPath := filepath.Clean(path)
	cleanDest := filepath.Clean(destDir)

	rel, err := filepath.Rel(cleanDest, cleanPath)
	if err != nil {
		return false
	}

	return !strings.HasPrefix(rel, "..") && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

// copyFile копирует содержимое файла из архива
func (u *UnzipService) copyFile(f *zip.File, destPath string) error {
	// Открываем исходный файл
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Создаем целевой файл
	dst, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer dst.Close()

	// Копируем содержимое
	_, err = io.Copy(dst, rc)
	return err
}
