package app

import (
	"converter/internal/config"
	"converter/internal/handler/convert"
	convertservice "converter/internal/service/convert"
	"log"
)

type App struct {
}

func NewApp() *App {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	_ = cfg

	im := convertservice.NewGoImageProcessor()
	md := convertservice.NewMDFinder(cfg.ConverterCfg.RootDir, ".md", im)
	c := convert.NewConvertHandler(md, im)

	c.HandleDirPipline(cfg.ConverterCfg.RootDir, cfg.ConverterCfg.TmpDir)
	return &App{}
}

func (a *App) Start() {

}
