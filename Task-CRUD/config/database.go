package config

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	db       *gorm.DB
	sqlDB    *sql.DB
	initOnce sync.Once
)

// InitPostgres menginisialisasi koneksi ke database PostgreSQL.
// Hanya akan dijalankan satu kali (singleton).
func InitPostgres(cfg *Config) (*gorm.DB, error) {
	var err error

	initOnce.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort,
		)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "public.",
				SingularTable: false,
			},
		})
		if err != nil {
			log.Printf("❌ Gagal menghubungkan ke PostgreSQL: %v", err)
			return
		}

		sqlDB, err = db.DB()
		if err != nil {
			log.Printf("❌ Gagal mendapatkan sql.DB dari GORM: %v", err)
			return
		}

		// Connection pooling
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
		sqlDB.SetConnMaxIdleTime(5 * time.Minute)

		log.Println("✅ Koneksi PostgreSQL berhasil dibuat dengan pooling")
	})

	// Jika db tetap nil, beri error manual
	if db == nil && err == nil {
		err = fmt.Errorf("koneksi database gagal tanpa error eksplisit")
	}

	return db, err
}

// GetDB mengembalikan instance *gorm.DB untuk digunakan di seluruh aplikasi.
func GetDB() *gorm.DB {
	if db == nil {
		log.Println("⚠️ DB belum diinisialisasi. Pastikan InitPostgres() sudah dipanggil.")
	}
	return db
}

// GetSqlDB mengembalikan instance *sql.DB jika perlu akses native SQL.
func GetSqlDB() *sql.DB {
	if sqlDB == nil {
		log.Println("⚠️ sqlDB belum diinisialisasi. Pastikan InitPostgres() sudah dipanggil.")
	}
	return sqlDB
}

// ClosePostgres menutup koneksi sqlDB secara aman saat shutdown aplikasi.
func ClosePostgres() error {
	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("⚠️ Gagal menutup koneksi database: %v", err)
			return err
		}
		log.Println("✅ Koneksi database berhasil ditutup")
	}
	return nil
}
