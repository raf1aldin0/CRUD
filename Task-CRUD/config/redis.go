package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis menginisialisasi koneksi Redis dengan konfigurasi dari Config
func InitRedis(cfg *Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword, // gunakan password jika ada
		DB:       0,                 // default DB
	})

	// Ping untuk cek koneksi
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		log.Printf("❌ Gagal terhubung ke Redis (%s:%s): %v", cfg.RedisHost, cfg.RedisPort, err)
		return err
	}

	log.Printf("✅ Berhasil terhubung ke Redis di %s:%s", cfg.RedisHost, cfg.RedisPort)
	return nil
}

// CloseRedis menutup koneksi Redis secara aman
func CloseRedis() error {
	if RedisClient == nil {
		return nil
	}
	log.Println("🚪 Menutup koneksi Redis...")
	return RedisClient.Close()
}
