package in

import (
	event "converter/internal/events"

	"github.com/IBM/sarama"
)

type Handler struct {
	eventHandler event.EventHandler
}

func NewHandler(eventHandler event.EventHandler) Handler {
	return Handler{
		eventHandler: eventHandler,
	}
}

func (h Handler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h Handler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		kafkaMsg := &event.KafkaMessage{
			Topic:     msg.Topic,
			Partition: msg.Partition,
			Offset:    msg.Offset,
			Key:       msg.Key,
			Value:     msg.Value,
			Headers:   convertHeaders(msg.Headers),
		}

		err := h.eventHandler.Handle(kafkaMsg)
		if err != nil {
			continue
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func convertHeaders(headers []*sarama.RecordHeader) map[string]string {
	result := make(map[string]string)
	for _, h := range headers {
		result[string(h.Key)] = string(h.Value)
	}
	return result
}
