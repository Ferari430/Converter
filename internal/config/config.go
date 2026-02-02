package config

import (
	"log"
	"runtime"
)

type Config struct {
	AppCfg       AppConfig
	ConverterCfg ConverterConfig
}

type AppConfig struct {
	Root           string
	Sep            string
	PandocPath     string
	WkhtmltopdfPdf string
}

type ConverterConfig struct {
	RootDir string // тут сканируем папку с файлами
	TmpDir  string // сюда сохраняем обработанные файлы
}

func LoadConfig() (*Config, error) {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err, "cant't load config file")
		return nil, err
	}

	return cfg, nil
}

func NewConfig() (*Config, error) {
	var (
		s           string
		pandoc      string
		wkhtmltopdf string
		r           string
		tmp         string
	)

	System := runtime.GOOS

	switch System {
	case "linux":
		s = `/`
		pandoc = "pandoc"
		wkhtmltopdf = "wkhtmltopdf"
		r = "/home/user/programmin/converter-20260131T211739Z-3-001/converter/testDir" //dir for converter
		tmp = "/home/user/programmin/converter-20260131T211739Z-3-001/converter/tmp"   // tmp dir
	case "windows":
		s = `\`
		pandoc = `C:\Program Files\Pandoc\pandoc.exe`
		wkhtmltopdf = `C:\Program Files\wkhtmltopdf\bin\wkhtmltopdf.exe`
		r = `B:\programmin-20260114T065921Z-1-001\programmin\converter\testDir`
		tmp = `B:\programmin-20260114T065921Z-1-001\programmin\converter\tmp`
	}

	return &Config{
		AppCfg: AppConfig{
			Sep:            s,
			PandocPath:     pandoc,
			WkhtmltopdfPdf: wkhtmltopdf,
		},
		ConverterCfg: ConverterConfig{
			RootDir: r,
			TmpDir:  tmp,
		},
	}, nil
}
