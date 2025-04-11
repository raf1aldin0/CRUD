package delivery

import (
	httpDelivery "Task-CRUD/delivery/http"
	repoRepo "Task-CRUD/internal/repository/repo"
	userRepo "Task-CRUD/internal/repository/user"
	"Task-CRUD/internal/usecase"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// NewRouter menerima *gorm.DB, *sql.DB, Redis client, dan Kafka writer
func NewRouter(gormDB *gorm.DB, sqlDB *sql.DB, rdb *redis.Client, kafkaWriter *kafka.Writer) *mux.Router {
	router := mux.NewRouter()

	// ===== Health Check =====
	router.HandleFunc("/health/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "live"})
	}).Methods("GET")

	router.HandleFunc("/health/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := sqlDB.Ping(); err != nil {
			log.Println("❌ SQL DB not ready:", err)
			http.Error(w, `{"status":"SQL DB not ready"}`, http.StatusServiceUnavailable)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Println("❌ Redis not ready:", err)
			http.Error(w, `{"status":"Redis not ready"}`, http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	}).Methods("GET")

	// ===== Dependency Injection =====

	// User (pakai SQL native dan Redis)
	userRepository := userRepo.NewUserRepositoryPostgres(sqlDB)
	userUseCase := usecase.NewUserUseCaseWithCache(userRepository, rdb)
	userHandler := httpDelivery.NewUserHandler(userUseCase)

	// Repository (pakai GORM + Redis + Kafka + Circuit Breaker + Tracing)
	repoRepository := repoRepo.NewRepoRepositoryGorm(gormDB)
	repoUseCase := usecase.NewRepoUseCaseFull(repoRepository, rdb, kafkaWriter)
	repoHandler := httpDelivery.NewRepoHandler(repoUseCase)

	// ===== User Routes =====
	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", userHandler.GetUsers).Methods("GET")
	userRouter.HandleFunc("/{id}", userHandler.GetUserByID).Methods("GET")
	userRouter.HandleFunc("", userHandler.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/{id}", userHandler.DeleteUser).Methods("DELETE")

	// ===== Repository Routes =====
	repoRouter := router.PathPrefix("/repositories").Subrouter()
	repoRouter.HandleFunc("", repoHandler.GetAllRepos).Methods("GET")
	repoRouter.HandleFunc("/{id}", repoHandler.GetRepositoryByID).Methods("GET")
	repoRouter.HandleFunc("", repoHandler.CreateRepo).Methods("POST")
	repoRouter.HandleFunc("/{id}", repoHandler.UpdateRepo).Methods("PUT")
	repoRouter.HandleFunc("/{id}", repoHandler.DeleteRepo).Methods("DELETE")

	return router
}
