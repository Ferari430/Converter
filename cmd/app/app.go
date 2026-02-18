package app

import (
	"context"
	kafka_in "converter/internal/adapters/in/kafka"
	kafka_out "converter/internal/adapters/out/kafka"
	"converter/internal/config"
	infra_kafka "converter/internal/infra/kafka"
	"converter/internal/infra/repository"
	"converter/internal/usecase"
	convertuc "converter/internal/usecase/convert"
	unzipuc "converter/internal/usecase/unzip"
	"log"
)

type App struct {
	consumer    infra_kafka.Consumer
	config      *config.Config
	handler     *kafka_in.KafkaMessageHandler
	kafkaClient *infra_kafka.KafkaClient
}

func NewApp() *App {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	kafkaClient, err := infra_kafka.NewClient(cfg.KafkaCfg)
	if err != nil {
		log.Fatal("Failed to create Kafka client:", err)
	}

	// Инициализируем usecase сервисы
	unzipService := unzipuc.NewUnzipService()
	imageProcessor := convertuc.NewGoImageProcessor()
	mdFinder := convertuc.NewMarkdownFinder(cfg.ConverterCfg.RootDir, ".md", imageProcessor)

	unzipRepo := repository.NewUnzipRepository(unzipService)
	convertRepo := repository.NewConvertRepository(mdFinder, imageProcessor, &cfg.ConverterCfg)

	producer, err := kafkaClient.Producer()
	if err != nil {
		log.Fatal("Failed to create producer:", err)
	}

	eventPublisher := kafka_out.NewKafkaEventPublisher(producer, cfg.KafkaCfg.Topic)

	convertUseCase := usecase.NewConvertArchiveUseCase(unzipRepo, convertRepo, eventPublisher)

	messageHandler := kafka_in.NewKafkaMessageHandler(convertUseCase)

	consumer, err := infra_kafka.NewConsumer(kafkaClient, cfg.KafkaCfg.ConsumerGroupID, messageHandler, cfg.KafkaCfg.Topic)
	if err != nil {
		log.Fatal("Failed to create consumer:", err)
	}

	return &App{
		consumer:    consumer,
		config:      cfg,
		handler:     messageHandler,
		kafkaClient: kafkaClient,
	}
}

func (a *App) Start(ctx context.Context) error {
	log.Println("Starting converter application...")
	log.Printf("Config: BrokersAddr=%s, ConsumerGroupID=%s, Topic=%s",
		a.config.KafkaCfg.BrokersAddr,
		a.config.KafkaCfg.ConsumerGroupID,
		a.config.KafkaCfg.Topic,
	)

	go func() {
		if err := a.consumer.Consume(ctx); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	log.Println("Converter is running. Press Ctrl+C to exit.")
	return nil
}

func (a *App) Stop() error {
	log.Println("Stopping converter application...")
	if a.consumer != nil {
		if err := a.consumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v", err)
		}
	}
	if a.kafkaClient != nil {
		if err := a.kafkaClient.Close(); err != nil {
			log.Printf("Error closing Kafka client: %v", err)
		}
	}
	return nil
}
