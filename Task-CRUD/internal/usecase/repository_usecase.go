package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"Task-CRUD/internal/entity"
	"Task-CRUD/internal/repository/repo" // ⬅️ import langsung ke package `repo`

	"github.com/redis/go-redis/v9"
)

type RepoUseCaseInterface interface {
	GetAllRepos(ctx context.Context) ([]entity.Repository, error)
	GetRepositoryByID(ctx context.Context, id uint) (*entity.Repository, error)
	CreateRepo(ctx context.Context, repo *entity.Repository) error
	UpdateRepo(ctx context.Context, id uint, repo *entity.Repository) error
	DeleteRepo(ctx context.Context, id uint) error
}

type RepoUseCase struct {
	repoRepo repo.RepoRepositoryInterfaceGorm // ⬅️ gunakan interface dari package repo langsung
	redis    *redis.Client
}

// Constructor tanpa Redis
func NewRepoUseCase(repoRepo repo.RepoRepositoryInterfaceGorm) RepoUseCaseInterface {
	return &RepoUseCase{
		repoRepo: repoRepo,
	}
}

// Constructor dengan Redis
func NewRepoUseCaseWithCache(repoRepo repo.RepoRepositoryInterfaceGorm, redisClient *redis.Client) RepoUseCaseInterface {
	return &RepoUseCase{
		repoRepo: repoRepo,
		redis:    redisClient,
	}
}

func (uc *RepoUseCase) GetAllRepos(ctx context.Context) ([]entity.Repository, error) {
	cacheKey := "repositories:all"

	if uc.redis != nil {
		cached, err := uc.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedRepos []entity.Repository
			if err := json.Unmarshal([]byte(cached), &cachedRepos); err == nil {
				fmt.Println("✅ Data repositories diambil dari Redis")
				return cachedRepos, nil
			}
			fmt.Printf("⚠️ Gagal unmarshal data Redis: %v\n", err)
		} else if err != redis.Nil {
			fmt.Printf("⚠️ Gagal ambil cache dari Redis: %v\n", err)
		}
	}

	repos, err := uc.repoRepo.GetAllRepositories()
	if err != nil {
		return nil, err
	}

	if uc.redis != nil {
		bytes, _ := json.Marshal(repos)
		if err := uc.redis.Set(ctx, cacheKey, bytes, 10*time.Minute).Err(); err != nil {
			fmt.Printf("⚠️ Gagal menyimpan cache ke Redis: %v\n", err)
		}
	}

	return repos, nil
}

func (uc *RepoUseCase) GetRepositoryByID(ctx context.Context, id uint) (*entity.Repository, error) {
	cacheKey := fmt.Sprintf("repository:%d", id)

	if uc.redis != nil {
		cached, err := uc.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedRepo entity.Repository
			if err := json.Unmarshal([]byte(cached), &cachedRepo); err == nil {
				fmt.Println("✅ Repository ditemukan di Redis")
				return &cachedRepo, nil
			}
			fmt.Printf("⚠️ Gagal unmarshal repository dari Redis: %v\n", err)
		} else if err != redis.Nil {
			fmt.Printf("⚠️ Redis error saat get: %v\n", err)
		}
	}

	// Ambil dari DB kalau tidak ada di Redis
	repo, err := uc.repoRepo.GetRepositoryByID(id)
	if err != nil {
		return nil, err
	}

	// Simpan ke Redis
	if uc.redis != nil {
		bytes, _ := json.Marshal(repo)
		if err := uc.redis.Set(ctx, cacheKey, bytes, 10*time.Minute).Err(); err != nil {
			fmt.Printf("⚠️ Gagal simpan repository ke Redis: %v\n", err)
		}
	}

	return repo, nil
}

func (uc *RepoUseCase) CreateRepo(ctx context.Context, repo *entity.Repository) error {
	if repo.Name == "" {
		return errors.New("nama repository tidak boleh kosong")
	}

	if err := uc.repoRepo.CreateRepository(repo); err != nil {
		return err
	}

	if uc.redis != nil {
		if err := uc.redis.Del(ctx, "repositories:all").Err(); err != nil {
			fmt.Printf("⚠️ Gagal hapus cache Redis setelah Create: %v\n", err)
		}
	}

	return nil
}

func (uc *RepoUseCase) UpdateRepo(ctx context.Context, id uint, repo *entity.Repository) error {
	if repo.Name == "" {
		return errors.New("nama repository tidak boleh kosong")
	}

	if err := uc.repoRepo.UpdateRepository(id, repo); err != nil {
		return err
	}

	if uc.redis != nil {
		if err := uc.redis.Del(ctx, "repositories:all").Err(); err != nil {
			fmt.Printf("⚠️ Gagal hapus cache Redis setelah Update: %v\n", err)
		}
	}

	return nil
}

func (uc *RepoUseCase) DeleteRepo(ctx context.Context, id uint) error {
	if err := uc.repoRepo.DeleteRepository(id); err != nil {
		return err
	}

	if uc.redis != nil {
		if err := uc.redis.Del(ctx, "repositories:all").Err(); err != nil {
			fmt.Printf("⚠️ Gagal hapus cache Redis setelah Delete: %v\n", err)
		}
	}

	return nil
}
