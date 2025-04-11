package http

import (
	"context"
	"net/http"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthHandler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewHealthHandler(db *gorm.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		DB:    db,
		Redis: redis,
	}
}

func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("👍 Liveness OK"))
}

func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	sqlDB, err := h.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		http.Error(w, "❌ DB not ready", http.StatusServiceUnavailable)
		return
	}

	if h.Redis != nil {
		if err := h.Redis.Ping(context.Background()).Err(); err != nil {
			http.Error(w, "❌ Redis not ready", http.StatusServiceUnavailable)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ Readiness OK"))
}
