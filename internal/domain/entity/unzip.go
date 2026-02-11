package entity

// UnzipResult - результат распаковки архива
type UnzipResult struct {
	DestDir string   // директория куда распаковали
	Files   []string // список файлов которые распаковали
	Count   int      // количество файлов
}
