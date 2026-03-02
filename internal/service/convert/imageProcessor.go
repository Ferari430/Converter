package convertservice

import (
	"converter/internal/repo"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type GoImageProcessor struct {
	db *repo.InMemoryDatabase
}

func NewGoImageProcessor() *GoImageProcessor {
	return &GoImageProcessor{
		db: nil, // Будет установлена через SetDB или передана при создании
	}
}

func NewGoImageProcessorWithDB(db *repo.InMemoryDatabase) *GoImageProcessor {
	return &GoImageProcessor{
		db: db,
	}
}

func (g *GoImageProcessor) RecursiveFindImage(root string) error {
	if g.db == nil {
		log.Println("[GoImageProcessor] ERROR: Database not initialized")
		return nil
	}

	extensions := []string{".png", ".jpg", ".jpeg"}
	log.Println("searching for images recursively in ", root)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Пропускаем ошибки доступа
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		for _, imageExt := range extensions {
			if ext == imageExt {
				filename := filepath.Base(path)
				g.db.AddImage(filename, path)
			}
		}
		return nil
	})
	return err
}

// todo: обработать пустую мапу или че то типо того
func (g *GoImageProcessor) GetInfo() map[string]string {
	if g.db == nil {
		return nil
	}
	return g.db.GetImages()
}
