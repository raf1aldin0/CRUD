package main

import (
	"Task-CRUD/config"
	"Task-CRUD/delivery"
	"Task-CRUD/internal/entity"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("📦 Memulai inisialisasi server...")

	// Load konfigurasi dari .env
	cfg := config.LoadConfig()
	log.Println("🔧 Konfigurasi berhasil dimuat")

	// Validasi konfigurasi penting
	if cfg.ServerPort == "" || cfg.DbName == "" || cfg.DbHost == "" || cfg.HttpReadTimeout == 0 {
		log.Fatal("❌ Konfigurasi tidak lengkap atau nilai timeout tidak di-set. Mohon cek file .env kamu")
	}

	// Inisialisasi PostgreSQL
	db, err := config.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("❌ Gagal inisialisasi PostgreSQL: %v", err)
	}
	log.Println("✅ Koneksi ke PostgreSQL berhasil")

	// AutoMigrate untuk entity
	if err := db.AutoMigrate(&entity.User{}, &entity.Repository{}); err != nil {
		log.Fatalf("❌ Gagal AutoMigrate: %v", err)
	}
	log.Println("✅ AutoMigrate berhasil")

	// Inisialisasi Redis
	if err := config.InitRedis(cfg); err != nil {
		log.Fatalf("❌ Gagal menginisialisasi Redis: %v", err)
	}
	log.Println("✅ Redis berhasil terhubung")

	// Setup router
	router := delivery.NewRouter(db, config.RedisClient)

	// Setup HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  cfg.HttpReadTimeout,
		WriteTimeout: cfg.HttpWriteTimeout,
		IdleTimeout:  cfg.HttpIdleTimeout,
	}

	// Jalankan server dalam goroutine
	go func() {
		log.Printf("🚀 Server berjalan di port %s...", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	<-shutdownChan
	log.Println("🛑 Mematikan server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("❌ Gagal shutdown server dengan baik: %v", err)
	}

	safeClose("PostgreSQL", config.ClosePostgres)
	safeClose("Redis", config.CloseRedis)

	log.Println("👋 Server dimatikan dengan aman")
}

// Helper untuk menutup resource dengan log
func safeClose(name string, closer func() error) {
	if err := closer(); err != nil {
		log.Printf("⚠️ Gagal menutup %s: %v", name, err)
	} else {
		log.Printf("✅ %s berhasil ditutup", name)
	}
}
