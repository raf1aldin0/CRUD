package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort       string
	DbHost           string
	DbPort           string
	DbUser           string
	DbPassword       string
	DbName           string
	DbMaxOpenConns   int
	DbMaxIdleConns   int
	DbConnMaxLifeSec int

	RedisHost     string
	RedisPort     string
	RedisPassword string

	KafkaBroker string
	KafkaTopic  string

	HttpReadTimeout  time.Duration
	HttpWriteTimeout time.Duration
	HttpIdleTimeout  time.Duration
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️  .env file tidak ditemukan, gunakan nilai default/env dari OS...")
	}

	// Default values
	viper.SetDefault("SERVER_PORT", "8080")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "postgres")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 20)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", 300)

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "secret123")

	viper.SetDefault("KAFKA_BROKER", "kafka:9092")
	viper.SetDefault("KAFKA_TOPIC", "repository-topic")

	viper.SetDefault("HTTP_READ_TIMEOUT", 15)
	viper.SetDefault("HTTP_WRITE_TIMEOUT", 15)
	viper.SetDefault("HTTP_IDLE_TIMEOUT", 60)

	cfg := &Config{
		ServerPort:       viper.GetString("SERVER_PORT"),
		DbHost:           viper.GetString("DB_HOST"),
		DbPort:           viper.GetString("DB_PORT"),
		DbUser:           viper.GetString("DB_USER"),
		DbPassword:       viper.GetString("DB_PASSWORD"),
		DbName:           viper.GetString("DB_NAME"),
		DbMaxOpenConns:   viper.GetInt("DB_MAX_OPEN_CONNS"),
		DbMaxIdleConns:   viper.GetInt("DB_MAX_IDLE_CONNS"),
		DbConnMaxLifeSec: viper.GetInt("DB_CONN_MAX_LIFETIME"),
		RedisHost:        viper.GetString("REDIS_HOST"),
		RedisPort:        viper.GetString("REDIS_PORT"),
		RedisPassword:    viper.GetString("REDIS_PASSWORD"),
		KafkaBroker:      viper.GetString("KAFKA_BROKER"),
		KafkaTopic:       viper.GetString("KAFKA_TOPIC"),
		HttpReadTimeout:  time.Duration(viper.GetInt("HTTP_READ_TIMEOUT")) * time.Second,
		HttpWriteTimeout: time.Duration(viper.GetInt("HTTP_WRITE_TIMEOUT")) * time.Second,
		HttpIdleTimeout:  time.Duration(viper.GetInt("HTTP_IDLE_TIMEOUT")) * time.Second,
	}

	// Validasi
	if cfg.ServerPort == "" || cfg.DbHost == "" || cfg.DbUser == "" || cfg.DbName == "" {
		log.Fatal("❌ Konfigurasi server/PostgreSQL tidak lengkap")
	}
	if cfg.RedisHost == "" || cfg.RedisPort == "" {
		log.Fatal("❌ Konfigurasi Redis tidak lengkap")
	}
	if cfg.KafkaBroker == "" || cfg.KafkaTopic == "" {
		log.Fatal("❌ Konfigurasi Kafka tidak lengkap")
	}

	log.Println("✅ Konfigurasi berhasil dimuat")
	return cfg
}
