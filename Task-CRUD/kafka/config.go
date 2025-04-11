package kafka

import (
	"os"
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

// LoadKafkaConfig membaca konfigurasi Kafka dari environment variable
func LoadKafkaConfig() KafkaConfig {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092" // Default jika tidak di-set
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "repository-topic" // Default jika tidak di-set
	}

	return KafkaConfig{
		Brokers: []string{broker},
		Topic:   topic,
	}
}
