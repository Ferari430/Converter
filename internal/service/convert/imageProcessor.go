package convertservice

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type GoImageProcessor struct {
	images map[string]string // map{pngName:pngPath}
}

func NewGoImageProcessor() *GoImageProcessor {
	m := make(map[string]string)
	return &GoImageProcessor{images: m}
}

func (g *GoImageProcessor) RecursiveFindImage(root string) error {
	extensions := []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".svg", ".webp"}
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

				if _, exists := g.images[filename]; !exists {
					g.images[filename] = path
				}
			}
		}
		return nil
	})
	return err
}

// todo: обработать пустую мапу или че то типо того
func (g *GoImageProcessor) GetInfo() map[string]string {
	if len(g.images) == 0 {
		return nil
	}
	
	return g.images
}
