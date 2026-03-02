package app

import (
	"context"
	"converter/internal/adapters/in"
	"converter/internal/adapters/out"
	"converter/internal/config"
	"converter/internal/handler/convert"
	"converter/internal/handler/pipeline"
	"converter/internal/handler/unzip"
	kafka2 "converter/internal/infra/kafka"
	"converter/internal/repo"
	convertservice "converter/internal/service/convert"
	unzip2 "converter/internal/service/unzip"
	"log"
)

type App struct {
	ConvertHandler *convert.ConvertHandler
	UnzipHandler   *unzip.UnzipHandler
	cfg            *config.Config
	consumer       in.Consumer
	db             *repo.InMemoryDatabase
}

func NewApp() *App {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	// Создаём централизованную БД
	db := repo.NewInMemoryDatabase()

	// Создаём сервисы с передачей БД
	im := convertservice.NewGoImageProcessorWithDB(db)
	md := convertservice.NewMDFinderWithDBAndWkhtmltopdf(cfg.ConverterCfg.RootDir, ".md", im, db, cfg.AppCfg.WkhtmltopdfPdf)
	c := convert.NewConvertHandler(md, im)

	///UNZIP

	s := unzip2.NewUnzipService()
	h := unzip.NewUnzipHandler(s)

	kafka, err := kafka2.NewClient(cfg.KafkaCfg)
	if err != nil {
		log.Fatal(err)
	}

	prod, err := out.NewProducer(kafka, cfg.KafkaCfg.Topic)
	_ = prod

	// Создаём архивный pipeline обработки
	processingPipeline := pipeline.NewArchiveProcessingPipeline(h, c, cfg.ConverterCfg.TmpDir)

	// Создаём обработчик Kafka событий
	kafkaEventHandler := in.NewKafkaEventHandler(processingPipeline)
	kafkaHandler := in.NewHandler(kafkaEventHandler)
	cons, err := in.NewConsumer(kafka, cfg.KafkaCfg.ConsumerGroupID, kafkaHandler)

	return &App{
		ConvertHandler: c,
		cfg:            cfg,
		UnzipHandler:   h,
		consumer:       cons,
		db:             db,
	}
}

// kafka
func (a *App) Start() {
	go func() {
		err := a.consumer.Consume(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("starting consumer")
	select {} // Ждём вечно, консьюмер обработает сообщения
}
