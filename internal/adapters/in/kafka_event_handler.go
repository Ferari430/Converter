package in

import (
	event "converter/internal/events"
	"converter/internal/handler/pipeline"
	"encoding/json"
	"log"
)

// KafkaEventHandler парсит события из Kafka
type KafkaEventHandler struct {
	pipeline *pipeline.ArchiveProcessingPipeline
}

func NewKafkaEventHandler(p *pipeline.ArchiveProcessingPipeline) *KafkaEventHandler {
	return &KafkaEventHandler{
		pipeline: p,
	}
}

// Handle парсит Kafka сообщение и передаёт в pipeline обработки
func (keh *KafkaEventHandler) Handle(msg *event.KafkaMessage) error {
	log.Println("[KafkaEventHandler] Received message", string(msg.Key), string(msg.Value))

	obj := make(map[string]interface{})
	err := json.Unmarshal(msg.Value, &obj)
	if err != nil {
		log.Println("[KafkaEventHandler] Error unmarshaling message:", err)
		return err
	}

	log.Println("[KafkaEventHandler] Parsed event:", obj)

	filePath, ok := obj["file_path"]
	if !ok {
		log.Println("[KafkaEventHandler] File path not found in event")
		return nil
	}

	archivePath := filePath.(string)
	log.Println("[KafkaEventHandler] Processing archive at:", archivePath)

	// Передаём архив в pipeline обработки
	err = keh.pipeline.Process(archivePath)
	if err != nil {
		log.Println("[KafkaEventHandler] Error processing archive:", err)
		return err
	}

	return nil
}
