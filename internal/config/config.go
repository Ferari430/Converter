package config

import (
	"flag"
	"os"
	"runtime"
	"strconv"

	"log"

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
		vps         bool
		envPath     string
	)

	flag.BoolVar(&vps, "vps", true, "is this running on vps?")
	flag.Parse()
	value, ok := os.LookupEnv("vps")
	if ok {
		boolVal, _ := strconv.ParseBool(value)
		vps = boolVal
	}

	System := runtime.GOOS
	if vps {
		envPath = "root/server/Converter/.env"
		log.Println("vps path for env file:", envPath)
	} else {
		envPath = `/home/user/programmin/obsidian_Project/prog/converter/.env`
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Println("Error loading .env file")
		return nil, err
	}

	switch System {
	case "linux":
		s = `/`
		pandoc = "pandoc"
		wkhtmltopdf = "wkhtmltopdf"
		r = "/home/user/programmin/converter-20260131T211739Z-3-001/converter/testDir" //dir for converter
		tmp = "/home/user/data/converter"                                              // tmp dir  сюда сохраняются обработанные файлы
	case "windows":
		s = `\`
		pandoc = `C:\Program Files\Pandoc\pandoc.exe`
		wkhtmltopdf = `C:\Program Files\wkhtmltopdf\bin\wkhtmltopdf.exe`
		r = `B:\programmin-20260114T065921Z-1-001\programmin\converter\test`                        //dir for converter
		tmp = `B:\programmin-20260114T065921Z-1-001\programmin\obsidian_Project\prog\converter\tmp` // tmp dir
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
