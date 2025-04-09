package delivery

import (
	httpDelivery "Task-CRUD/delivery/http"
	"Task-CRUD/internal/repository/repo"
	"Task-CRUD/internal/repository/user"
	"Task-CRUD/internal/usecase"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, rdb *redis.Client) *mux.Router {
	router := mux.NewRouter()

	// ========== Health Check Handlers ==========

	// Liveness Probe
	router.HandleFunc("/health/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "live"})
	}).Methods("GET")

	// Readiness Probe
	router.HandleFunc("/health/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check DB connection
		sqlDB, err := db.DB()
		if err != nil {
			log.Println("❌ DB connection error:", err)
			http.Error(w, `{"status":"DB connection error"}`, http.StatusServiceUnavailable)
			return
		}
		if err := sqlDB.Ping(); err != nil {
			log.Println("❌ DB not ready:", err)
			http.Error(w, `{"status":"DB not ready"}`, http.StatusServiceUnavailable)
			return
		}

		// Check Redis connection
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Println("❌ Redis not ready:", err)
			http.Error(w, `{"status":"Redis not ready"}`, http.StatusServiceUnavailable)
			return
		}

		// All good
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	}).Methods("GET")

	// ========== Dependency Injection ==========

	userRepo := user.NewUserRepositoryGorm(db)
	userUC := usecase.NewUserUseCaseWithCache(userRepo, rdb)
	userHandler := httpDelivery.NewUserHandler(userUC)

	repoRepo := repo.NewRepoRepositoryGorm(db)
	repoUC := usecase.NewRepoUseCaseWithCache(repoRepo, rdb)
	repoHandler := httpDelivery.NewRepoHandler(repoUC)

	// ========== User Routes ==========
	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", userHandler.GetUsers).Methods("GET")
	userRouter.HandleFunc("/{id}", userHandler.GetUserByID).Methods("GET")
	userRouter.HandleFunc("", userHandler.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/{id}", userHandler.DeleteUser).Methods("DELETE")

	// ========== Repository Routes ==========
	repoRouter := router.PathPrefix("/repositories").Subrouter()
	repoRouter.HandleFunc("", repoHandler.GetAllRepos).Methods("GET")
	repoRouter.HandleFunc("/{id}", repoHandler.GetRepositoryByID).Methods("GET")
	repoRouter.HandleFunc("", repoHandler.CreateRepo).Methods("POST")
	repoRouter.HandleFunc("/{id}", repoHandler.UpdateRepo).Methods("PUT")
	repoRouter.HandleFunc("/{id}", repoHandler.DeleteRepo).Methods("DELETE")

	return router
}
