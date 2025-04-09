package delivery

import (
	httpDelivery "Task-CRUD/delivery/http"
	"Task-CRUD/internal/repository/repo"
	"Task-CRUD/internal/repository/user"
	"Task-CRUD/internal/usecase"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, rdb *redis.Client) *mux.Router {
	router := mux.NewRouter()

	// ========== Probes ==========
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"live"}`))
	}).Methods("GET")

	router.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB()
		if err != nil {
			log.Println("❌ Gagal mendapatkan koneksi DB:", err)
			http.Error(w, `{"status":"DB connection error"}`, http.StatusServiceUnavailable)
			return
		}
		if err := sqlDB.Ping(); err != nil {
			log.Println("❌ Database belum siap:", err)
			http.Error(w, `{"status":"DB not ready"}`, http.StatusServiceUnavailable)
			return
		}

		// Redis Check
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if _, err := rdb.Ping(ctx).Result(); err != nil {
			log.Println("❌ Redis belum siap:", err)
			http.Error(w, `{"status":"Redis not ready"}`, http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
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
