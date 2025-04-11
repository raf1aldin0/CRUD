package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(cfg KafkaConfig) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		GroupID: "task-crud-group",
	})

	go func() {
		fmt.Println("🟢 Kafka consumer listening...")
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("❌ Error reading Kafka message: %v", err)
				continue
			}
			fmt.Printf("📥 Kafka received: %s\n", string(m.Value))
		}
	}()
}
