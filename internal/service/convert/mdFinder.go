package convertservice

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

type imageFinder interface {
	GetInfo() map[string]string
}

type MDFinder struct {
	root        string
	dataType    string
	Data        map[string]string
	ImageFinder imageFinder
}

type Config struct {
	inputPath  string
	outputPath string
}

func NewMDFinder(root string, t string, imageFinder imageFinder) *MDFinder {
	return &MDFinder{
		root:        root,
		dataType:    t,
		Data:        make(map[string]string),
		ImageFinder: imageFinder,
	}
}

func (f *MDFinder) ScanFiles(root string) (map[string]string, error) {

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
		} else {

		}

		return nil
	})
	log.Println(a)
	return f.Data, err
}

func (f *MDFinder) ProcessFile(inputPath, tmpdir string) (string, error) {
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

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		processedLine := processLine(line, f.ImageFinder.GetInfo(), tmpdir)
		_, err = writer.WriteString(processedLine + "\n")
		if err != nil {
			log.Println(err)
		}
	}

	return outputPath, scanner.Err()
}

var (
	// Основное регулярное выражение - ловим ВСЁ внутри скобок
	imageRegex = regexp.MustCompile(`!\[(\[)?([^\]]+?)(\])?\]`)

	// Дополнительное: для извлечения чистого имени файла (без пути в скобках)
	fullLinkRegex = regexp.MustCompile(`!\[(\[)?([^\]]+?)(\])?\]\(([^)]+)\)`)
)

func processLine(line string, imagePath map[string]string, output string) string {
	if strings.Contains(line, "![") {
		return handleLineWithImage(line, imagePath, output)
	}
	return line
}

func handleLineWithImage(l string, imagePath map[string]string, output string) string {
	// Сначала обрабатываем ссылки с путями в скобках
	result := fullLinkRegex.ReplaceAllStringFunc(l, func(match string) string {
		parts := fullLinkRegex.FindStringSubmatch(match)
		if len(parts) < 5 {
			return match
		}

		filename := cleanFilename(parts[2]) // Очищаем имя файла
		currentPath := strings.TrimSpace(parts[4])
		isDoubleBracket := parts[1] == "[" && parts[3] == "]"

		// Если путь уже есть - возможно, его не нужно менять
		if isGoodPath(currentPath) {
			return match
		}

		// Ищем файл в мапе
		return processImageLink(filename, currentPath, isDoubleBracket, imagePath, output)
	})

	// Затем обрабатываем простые ![[filename]]
	result = imageRegex.ReplaceAllStringFunc(result, func(match string) string {
		// Пропускаем, если это уже обработано как полная ссылка
		if strings.Contains(match, "](") {
			return match
		}

		parts := imageRegex.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}

		filename := cleanFilename(parts[2])
		isDoubleBracket := parts[1] == "[" && parts[3] == "]"

		// Ищем файл в мапе
		return processImageLink(filename, "", isDoubleBracket, imagePath, output)
	})

	return result
}

// Очистка имени файла от лишних символов
func cleanFilename(raw string) string {
	// Убираем пробелы по краям
	filename := strings.TrimSpace(raw)

	// Убираем markdown-экранирование (если есть)
	filename = strings.ReplaceAll(filename, `\ `, " ") // \_ -> _
	filename = strings.ReplaceAll(filename, `\[`, "[") // \[ -> [
	filename = strings.ReplaceAll(filename, `\]`, "]") // \] -> ]
	filename = strings.ReplaceAll(filename, `\(`, "(") // \( -> (
	filename = strings.ReplaceAll(filename, `\)`, ")") // \) -> )

	return filename
}

// Поиск файла с учетом пробелов и разных вариантов
func findImageInMap(filename string, imagePath map[string]string) (string, bool) {
	// Вариант 1: Точное совпадение
	if path, ok := imagePath[filename]; ok {
		return path, true
	}

	// Вариант 2: Без учета регистра
	lowerFilename := strings.ToLower(filename)
	for key, path := range imagePath {
		if strings.ToLower(key) == lowerFilename {
			return path, true
		}
	}

	// Вариант 3: Частичное совпадение (для файлов с путями)
	baseName := filepath.Base(filename)
	for key, path := range imagePath {
		if strings.ToLower(filepath.Base(key)) == strings.ToLower(baseName) {
			return path, true
		}
	}

	// Вариант 4: Ищем с добавлением расширений
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
	// Ищем файл
	path, found := findImageInMap(filename, imagePath)
	if !found && currentPath != "" {
		// Если не нашли в мапе, но есть текущий путь - используем его
		path = currentPath
		found = true
	}

	if !found {
		log.Printf("Не найден файл: %s", filename)
		// Возвращаем оригинал в правильном формате
		if isDoubleBracket {
			return fmt.Sprintf("![[%s]]", filename)
		}
		return fmt.Sprintf("![%s]", filename)
	}

	relativePath, err := filepath.Rel(output, path)
	if err != nil {
		log.Printf("Ошибка преобразования пути для %s: %v", filename, err)
		// Возвращаем оригинал
		if isDoubleBracket {
			return fmt.Sprintf("![[%s]]", filename)
		}
		return fmt.Sprintf("![%s]", filename)
	}

	// Форматируем результат
	relativePath = filepath.ToSlash(relativePath)

	// Экранируем пробелы в путях для markdown
	if strings.Contains(relativePath, " ") {
		relativePath = strings.ReplaceAll(relativePath, " ", "%20")
	}

	if isDoubleBracket {
		return fmt.Sprintf("![[%s]](%s)", filename, relativePath)
	}
	return fmt.Sprintf("![%s](%s)", filename, relativePath)
}

func isGoodPath(path string) bool {
	// Проверяем, является ли путь "хорошим" (относительным, рабочим)
	return strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../") ||
		strings.HasPrefix(path, "/") ||
		strings.Contains(path, "://") // http://, https://
}

func (f *MDFinder) ConvertMDToPDF(inputFile, outputFile, resourceBaseDir string) error {

	if _, err := os.Stat(inputFile); err != nil {
		return fmt.Errorf("input file not found: %s", inputFile)
	}

	PandocPath := `pandoc`

	// Преобразуем пути к абсолютным
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

	// Логируем для отладки
	log.Printf("Converting: %s -> %s", absInput, absOutput)
	log.Printf("Working dir: %s", inputDir)
	log.Printf("Resource dir: %s", resourceDir)

	// Формируем команду с новым синтаксисом
	cmd := exec.Command(PandocPath,
		absInput,
		"-o", absOutput,
		// Новый синтаксис вместо --self-contained
		"--embed-resources",
		"--standalone",
		// Настройки для игнорирования YAML
		//"--from=markdown",
		"--from=markdown-yaml_metadata_block",
		// Пути для поиска ресурсов
		"--resource-path="+resourceDir,
		"--verbose",
	)

	// Устанавливаем рабочую директорию
	cmd.Dir = inputDir

	// Запускаем команду
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pandoc error: %v\nOutput: %s", err, string(output))
	}

	log.Printf("Successfully converted %s to %s", absInput, absOutput)
	return nil
}
