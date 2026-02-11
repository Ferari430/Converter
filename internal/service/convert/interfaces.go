package convertservice

// MdFinder - интерфейс для поиска и обработки MD файлов
type MdFinder interface {
	ScanFiles(root string) (map[string]string, error)
	ProcessFile(inputPath, tmpdir string) (string, error)
	ConvertMDToPDF(inputFile, outputFile, resourceBaseDir string) error
}

// ImageFinder - интерфейс для поиска изображений
type ImageFinder interface {
	RecursiveFindImage(root string) error
	GetInfo() map[string]string
}
