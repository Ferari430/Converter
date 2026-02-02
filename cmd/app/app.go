package app

import (
	"converter/internal/config"
	"converter/internal/handler/convert"
	"converter/internal/handler/unzip"
	convertservice "converter/internal/service/convert"
	unzip2 "converter/internal/service/unzip"
	"log"
)

type App struct {
	ConvertHandler *convert.ConvertHandler
	UnzipHandler   *unzip.UnzipHandler
	cfg            *config.Config
	kafka          Kafka
}

type Kafka interface {
	Get() string
}

func NewApp() *App {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	im := convertservice.NewGoImageProcessor()
	md := convertservice.NewMDFinder(cfg.ConverterCfg.RootDir, ".md", im)
	c := convert.NewConvertHandler(md, im)

	///UNZIP

	s := unzip2.NewUnzipService()
	h := unzip.NewUnzipHandler(s)

	return &App{
		ConvertHandler: c,
		cfg:            cfg,
		UnzipHandler:   h,
	}
}

// kafka
func (a *App) Start() {
	//root := a.kafka.Get()
	//обработка сообщения
	//root := a.cfg.ConverterCfg.RootDir
	//log.Println("root:", root)
	//a.handler.HandleDirPipline(root, a.cfg.ConverterCfg.TmpDir)
	root := `B:\shadowarchive.zip`
	log.Println("start")
	result, err := a.UnzipHandler.Unzip(root)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(result)
}

func (a *App) Pipline() {
	//unzip  return filepath string

	//convert   return filepath string for PDFs files  or S3 link
}
