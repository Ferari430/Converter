package config

import (
	"log"
	"os"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	AppCfg       AppConfig
	ConverterCfg ConverterConfig
	KafkaCfg     KafkaConfig
}

type KafkaConfig struct {
	BrokersAddr     string
	ConsumerGroupID string
	Topic           string
}

type AppConfig struct {
	Root           string
	Sep            string
	PandocPath     string
	WkhtmltopdfPdf string
}

type ConverterConfig struct {
	RootDir string
	TmpDir  string
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

	err := godotenv.Load(`B:\programmin-20260114T065921Z-1-001\programmin\obsidian_Project\prog\converter\.env`)
	if err != nil {
		log.Println("Error loading .env file")
		return nil, err
	}

	switch System {
	case "linux":
		s = `/`
		pandoc = "pandoc"
		wkhtmltopdf = "wkhtmltopdf"
		r = "/home/user/programmin/converter-20260131T211739Z-3-001/converter/testDir"
		tmp = "/home/user/programmin/converter-20260131T211739Z-3-001/converter/tmp"
	case "windows":
		s = `\`
		pandoc = `C:\Program Files\Pandoc\pandoc.exe`
		wkhtmltopdf = `C:\Program Files\wkhtmltopdf\bin\wkhtmltopdf.exe`
		r = `B:\programmin-20260114T065921Z-1-001\programmin\converter\test`
		tmp = `B:\programmin-20260114T065921Z-1-001\programmin\obsidian_Project\prog\converter\tmp`
	}

	BrokersAddr := os.Getenv("KAFKA_PORT")
	ConsumerGroupID := os.Getenv("CONSUMER_GROUP_ID")
	Topic := os.Getenv("TOPIC")

	log.Println("BrokersAddr:", BrokersAddr, "ConsumerGroupID:", ConsumerGroupID, "Topic:", Topic)

	return &Config{
		AppCfg: AppConfig{
			Sep:            s,
			PandocPath:     pandoc,
			WkhtmltopdfPdf: wkhtmltopdf,
		},
		ConverterCfg: ConverterConfig{
			RootDir: r,
			TmpDir:  tmp,
		}, KafkaCfg: KafkaConfig{
			BrokersAddr:     BrokersAddr,
			ConsumerGroupID: ConsumerGroupID,
			Topic:           Topic,
		},
	}, nil
}
