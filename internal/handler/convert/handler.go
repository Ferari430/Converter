package convert

import (
	"converter/internal/config"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

type MdFinder interface {
	ScanFiles(root string) (map[string]string, error)
	ProcessFile(inputPath, tmpdir string) (string, error)
	ConvertMDToPDF(inputFile, outputFile, resourceBaseDir string) error
}

type ImageFinder interface {
	RecursiveFindImage(root string) error
	GetInfo() map[string]string
}

type ConvertHandler struct {
	md          MdFinder
	ch          chan string
	imageFinder ImageFinder
	cfg         *config.Config
}

func NewConvertHandler(mdFinder MdFinder, i ImageFinder) *ConvertHandler {
	return &ConvertHandler{
		md:          mdFinder,
		imageFinder: i,
	}
}

func (c *ConvertHandler) ScanMdFiles(root string) (map[string]string, error) {

	return c.md.ScanFiles(root)
}

func (c *ConvertHandler) ProcessFiles(files map[string]string, tmpdir string) ([]string, error) {
	tmpdirs := make([]string, 0, len(files))
	for _, root := range files {
		tmpFIle, err := c.md.ProcessFile(root, tmpdir)
		if err != nil {
			continue
		}

		tmpdirs = append(tmpdirs, tmpFIle)

	}

	return tmpdirs, nil
}

func (c *ConvertHandler) ConvertFiles(tmpFiles []string, tmpdir string) error {
	failedFiles := make([]string, 0, len(tmpFiles))
	var processedFiles int64
	source := filepath.Dir(tmpFiles[0])
	ok := false
	for _, file := range tmpFiles {
		log.Println("file:", file)
		htmlName := changeExtensionToHTML(file)
		log.Println(file, tmpdir)
		err := c.md.ConvertMDToPDF(file, htmlName, source)
		if err != nil {

			failedFiles = append(failedFiles, file)
			continue
		}
		processedFiles++
	}

	ok = true

	defer func() {
		if ok {
			log.Println("SUCCESS")
			fmt.Printf("Failed to convert files: %v\nConverted files: %v\nFilesCount:%v\n", len(failedFiles), processedFiles, len(tmpFiles))
			return
		}

		log.Println("FAIL")
		fmt.Printf("Failed to convert files: %v\nConverted files: %v\nFilesCount:%v\n", len(failedFiles), processedFiles, len(tmpFiles))
		return
	}()

	return nil
}

func (c *ConvertHandler) FindImages(root string) {

	err := c.imageFinder.RecursiveFindImage(root)
	if err != nil {
		log.Println(err)
	}

}

func changeExtensionToHTML(filepathStr string) string {
	ext := filepath.Ext(filepathStr)
	base := strings.TrimSuffix(filepathStr, ext)
	base = strings.TrimSuffix(base, "_TEMP")
	return base + ".html"
}
	ext := filepath.Ext(filepathStr)
	base := strings.TrimSuffix(filepathStr, ext)
	base = strings.TrimSuffix(base, "_TEMP")
	return base + ".html"
}

func (c *ConvertHandler) HandleDirPipline(root, tmpdir string) {

	log.Println("tmpdir:", tmpdir)
	c.FindImages(root)

	files, err := c.ScanMdFiles(root)

	if err != nil {
		log.Println(err)
	}

	tmpfiles, err := c.ProcessFiles(files, tmpdir)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.ConvertFiles(tmpfiles, tmpdir)

	if err != nil {
		log.Fatal(err)
	}
}
