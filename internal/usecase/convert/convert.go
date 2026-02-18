package convert

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// ImageFinder - интерфейс для поиска изображений
type ImageFinder interface {
	RecursiveFindImage(root string) error
	GetInfo() map[string]string
}

// MDFinder - интерфейс для поиска и обработки MD файлов
type MDFinder interface {
	ScanFiles(root string) (map[string]string, error)
	ProcessFile(inputPath, tmpdir string) (string, error)
	ConvertMDToPDF(inputFile, outputFile, resourceBaseDir string) error
}

// MarkdownFinder - реализация MDFinder
type MarkdownFinder struct {
	root        string
	dataType    string
	Data        map[string]string
	ImageFinder ImageFinder
}

// NewMarkdownFinder - конструктор для MarkdownFinder
func NewMarkdownFinder(root string, t string, imageFinder ImageFinder) *MarkdownFinder {
	return &MarkdownFinder{
		root:        root,
		dataType:    t,
		Data:        make(map[string]string),
		ImageFinder: imageFinder,
	}
}

// ScanFiles - сканирует директорию и ищет markdown файлы
func (f *MarkdownFinder) ScanFiles(root string) (map[string]string, error) {
	a := 0
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == f.dataType && !strings.Contains(path, "excalidraw") && !strings.Contains(path, ".convas") {
			filename := filepath.Base(path)
			log.Println(path, ext)
			if _, exists := f.Data[filename]; !exists {
				f.Data[filename] = path
				a++
			}
		}

		return nil
	})
	log.Println(a)
	return f.Data, err
}

// ProcessFile - обработает markdown файл и заменит ссылки на изображения
func (f *MarkdownFinder) ProcessFile(inputPath, tmpdir string) (string, error) {
	log.Printf("Processing file %s", inputPath)

	input, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}

	defer input.Close()

	if err := os.MkdirAll(tmpdir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать папку tmp в %s: %v", tmpdir, err)
	}

	inputFilename := filepath.Base(inputPath)

	outputPath := filepath.Join(tmpdir,
		fmt.Sprintf("%s_TEMP.md", strings.TrimSuffix(inputFilename, ".md")))
	output, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		processedLine := processLine(line, f.ImageFinder.GetInfo(), tmpdir)
		_, err = writer.WriteString(processedLine + "\n")
		if err != nil {
			log.Println(err)
		}
	}

	return outputPath, scanner.Err()
}

// ConvertMDToPDF - конвертирует markdown в HTML используя pandoc
func (f *MarkdownFinder) ConvertMDToPDF(inputFile, outputFile, resourceBaseDir string) error {
	if _, err := os.Stat(inputFile); err != nil {
		return fmt.Errorf("input file not found: %s", inputFile)
	}

	PandocPath := `pandoc`

	absInput, err := filepath.Abs(inputFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	absOutput, err := filepath.Abs(outputFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	inputDir := filepath.Dir(absInput)
	resourceDir := resourceBaseDir
	if resourceDir == "" {
		resourceDir = inputDir
	}

	log.Printf("Converting: %s -> %s", absInput, absOutput)
	log.Printf("Working dir: %s", inputDir)
	log.Printf("Resource dir: %s", resourceDir)

	cmd := exec.Command(PandocPath,
		absInput,
		"-o", absOutput,
		"--embed-resources",
		"--standalone",
		"--from=markdown-yaml_metadata_block",
		"--resource-path="+resourceDir,
		"--verbose",
	)

	cmd.Dir = inputDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pandoc error: %v\nOutput: %s", err, string(output))
	}

	log.Printf("Successfully converted %s to %s", absInput, absOutput)
	return nil
}

// GoImageProcessor - реализация ImageFinder
type GoImageProcessor struct {
	images map[string]string // map{imageName:imagePath}
}

// NewGoImageProcessor - конструктор для GoImageProcessor
func NewGoImageProcessor() *GoImageProcessor {
	return &GoImageProcessor{images: make(map[string]string)}
}

// RecursiveFindImage - рекурсивно ищет все изображения в директории
func (g *GoImageProcessor) RecursiveFindImage(root string) error {
	extensions := []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".svg", ".webp"}
	log.Println("searching for images recursively in ", root)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
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

// GetInfo - возвращает карту найденных изображений
func (g *GoImageProcessor) GetInfo() map[string]string {
	if len(g.images) == 0 {
		return nil
	}
	return g.images
}

// Helper functions

var (
	imageRegex    = regexp.MustCompile(`!\[(\[)?([^\]]+?)(\])?\]`)
	fullLinkRegex = regexp.MustCompile(`!\[(\[)?([^\]]+?)(\])?\]\(([^)]+)\)`)
)

func processLine(line string, imagePath map[string]string, output string) string {
	if strings.Contains(line, "---") {
		log.Println("processing invalid line")
		return ""
	}

	if strings.Contains(line, "![") {
		return handleLineWithImage(line, imagePath, output)
	}
	return line
}

func handleLineWithImage(l string, imagePath map[string]string, output string) string {
	result := fullLinkRegex.ReplaceAllStringFunc(l, func(match string) string {
		parts := fullLinkRegex.FindStringSubmatch(match)
		if len(parts) < 5 {
			return match
		}

		filename := cleanFilename(parts[2])
		currentPath := strings.TrimSpace(parts[4])
		isDoubleBracket := parts[1] == "[" && parts[3] == "]"

		if isGoodPath(currentPath) {
			return match
		}

		return processImageLink(filename, currentPath, isDoubleBracket, imagePath, output)
	})

	result = imageRegex.ReplaceAllStringFunc(result, func(match string) string {
		if strings.Contains(match, "](") {
			return match
		}

		parts := imageRegex.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}

		filename := cleanFilename(parts[2])
		isDoubleBracket := parts[1] == "[" && parts[3] == "]"

		return processImageLink(filename, "", isDoubleBracket, imagePath, output)
	})

	return result
}

func cleanFilename(raw string) string {
	filename := strings.TrimSpace(raw)

	filename = strings.ReplaceAll(filename, `\ `, " ")
	filename = strings.ReplaceAll(filename, `\[`, "[")
	filename = strings.ReplaceAll(filename, `\]`, "]")
	filename = strings.ReplaceAll(filename, `\(`, "(")
	filename = strings.ReplaceAll(filename, `\)`, ")")

	return filename
}

func findImageInMap(filename string, imagePath map[string]string) (string, bool) {
	if path, ok := imagePath[filename]; ok {
		return path, true
	}

	lowerFilename := strings.ToLower(filename)
	for key, path := range imagePath {
		if strings.ToLower(key) == lowerFilename {
			return path, true
		}
	}

	baseName := filepath.Base(filename)
	for key, path := range imagePath {
		if strings.EqualFold(filepath.Base(key), baseName) {
			return path, true
		}
	}

	if !strings.Contains(filename, ".") {
		extensions := []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".svg"}
		for _, ext := range extensions {
			if path, ok := imagePath[filename+ext]; ok {
				return path, true
			}
			if path, ok := imagePath[filename+strings.ToUpper(ext)]; ok {
				return path, true
			}
		}
	}

	return "", false
}

func processImageLink(filename, currentPath string, isDoubleBracket bool,
	imagePath map[string]string, output string) string {
	path, found := findImageInMap(filename, imagePath)
	if !found && currentPath != "" {
		path = currentPath
		found = true
	}

	if !found {
		log.Printf("Не найден файл: %s", filename)
		if isDoubleBracket {
			return fmt.Sprintf("![[%s]]", filename)
		}
		return fmt.Sprintf("![%s]", filename)
	}

	relativePath, err := filepath.Rel(output, path)
	if err != nil {
		log.Printf("Ошибка преобразования пути для %s: %v", filename, err)
		if isDoubleBracket {
			return fmt.Sprintf("![[%s]]", filename)
		}
		return fmt.Sprintf("![%s]", filename)
	}

	relativePath = filepath.ToSlash(relativePath)

	if strings.Contains(relativePath, " ") {
		relativePath = strings.ReplaceAll(relativePath, " ", "%20")
	}

	if isDoubleBracket {
		return fmt.Sprintf("![[%s]](%s)", filename, relativePath)
	}
	return fmt.Sprintf("![%s](%s)", filename, relativePath)
}

func isGoodPath(path string) bool {
	return strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../") ||
		strings.HasPrefix(path, "/") ||
		strings.Contains(path, "://")
}
