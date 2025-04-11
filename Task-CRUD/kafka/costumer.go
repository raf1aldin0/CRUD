package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
	Topic  string
}

func NewProducer(cfg KafkaConfig) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		Writer: writer,
		Topic:  cfg.Topic,
	}
}

func (kp *KafkaProducer) Publish(message interface{}) error {
	// Konversi message ke JSON
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Tambahkan context dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Kirim pesan ke Kafka
	err = kp.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("repository-event"),
		Value: data,
	})
	if err != nil {
		log.Printf("❌ Failed to publish Kafka message: %v", err)
	}
	return err
}
