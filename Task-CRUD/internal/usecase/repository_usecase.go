package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"Task-CRUD/internal/cbreaker"
	"Task-CRUD/internal/entity"
	interfaces "Task-CRUD/internal/interfaces"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
)

type RepoUseCase struct {
	repoRepo interfaces.RepoRepositoryInterfaceGorm
	redis    *redis.Client
	breaker  *gobreaker.CircuitBreaker
}

func NewRepoUseCase(repoRepo interfaces.RepoRepositoryInterfaceGorm) interfaces.RepoUseCaseInterface {
	return &RepoUseCase{
		repoRepo: repoRepo,
		breaker:  cbreaker.Breaker,
	}
}

func NewRepoUseCaseWithCache(repoRepo interfaces.RepoRepositoryInterfaceGorm, redisClient *redis.Client) interfaces.RepoUseCaseInterface {
	return &RepoUseCase{
		repoRepo: repoRepo,
		redis:    redisClient,
		breaker:  cbreaker.Breaker,
	}
}

func (uc *RepoUseCase) GetAllRepos(ctx context.Context) ([]entity.Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoUseCase.GetAllRepos")
	defer span.Finish()

	cacheKey := "repositories:all"

	if uc.redis != nil {
		cached, err := uc.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedRepos []entity.Repository
			if err := json.Unmarshal([]byte(cached), &cachedRepos); err == nil {
				span.LogFields(log.String("cache", "hit"))
				fmt.Println("✅ Data repositories diambil dari Redis")
				return cachedRepos, nil
			}
			span.LogFields(log.Error(err))
		} else if err != redis.Nil {
			span.LogFields(log.Error(err))
		}
	}

	result, err := uc.breaker.Execute(func() (interface{}, error) {
		return uc.repoRepo.GetAllRepositories(ctx)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return nil, err
	}
	repos := result.([]entity.Repository)

	if uc.redis != nil {
		bytes, _ := json.Marshal(repos)
		_ = uc.redis.Set(ctx, cacheKey, bytes, 10*time.Minute).Err()
	}

	return repos, nil
}

func (uc *RepoUseCase) GetRepositoryByID(ctx context.Context, id uint) (*entity.Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoUseCase.GetRepositoryByID")
	defer span.Finish()

	cacheKey := fmt.Sprintf("repository:%d", id)

	if uc.redis != nil {
		cached, err := uc.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedRepo entity.Repository
			if err := json.Unmarshal([]byte(cached), &cachedRepo); err == nil {
				span.LogFields(log.String("cache", "hit"))
				fmt.Println("✅ Repository ditemukan di Redis")
				return &cachedRepo, nil
			}
			span.LogFields(log.Error(err))
		}
	}

	result, err := uc.breaker.Execute(func() (interface{}, error) {
		return uc.repoRepo.GetRepositoryByID(ctx, id)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return nil, err
	}
	repo := result.(*entity.Repository)

	if uc.redis != nil {
		bytes, _ := json.Marshal(repo)
		_ = uc.redis.Set(ctx, cacheKey, bytes, 10*time.Minute).Err()
	}

	return repo, nil
}

func (uc *RepoUseCase) CreateRepo(ctx context.Context, repo *entity.Repository) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoUseCase.CreateRepo")
	defer span.Finish()

	if err := validateRepository(repo); err != nil {
		span.LogFields(log.Error(err))
		return err
	}

	_, err := uc.breaker.Execute(func() (interface{}, error) {
		return nil, uc.repoRepo.CreateRepository(ctx, repo)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return err
	}

	if uc.redis != nil {
		_ = uc.redis.Del(ctx, "repositories:all").Err()
	}

	return nil
}

func (uc *RepoUseCase) UpdateRepo(ctx context.Context, id uint, repo *entity.Repository) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoUseCase.UpdateRepo")
	defer span.Finish()

	if err := validateRepository(repo); err != nil {
		span.LogFields(log.Error(err))
		return err
	}

	_, err := uc.breaker.Execute(func() (interface{}, error) {
		return nil, uc.repoRepo.UpdateRepository(ctx, id, repo)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return err
	}

	if uc.redis != nil {
		_ = uc.redis.Del(ctx, "repositories:all").Err()
	}

	return nil
}

func (uc *RepoUseCase) DeleteRepo(ctx context.Context, id uint) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoUseCase.DeleteRepo")
	defer span.Finish()

	_, err := uc.breaker.Execute(func() (interface{}, error) {
		return nil, uc.repoRepo.DeleteRepository(ctx, id)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return err
	}

	if uc.redis != nil {
		_ = uc.redis.Del(ctx, "repositories:all").Err()
	}

	return nil
}

func validateRepository(repo *entity.Repository) error {
	if repo.Name == "" {
		return errors.New("nama repository tidak boleh kosong")
	}
	if repo.URL == "" {
		return errors.New("URL repository tidak boleh kosong")
	}
	if _, err := url.ParseRequestURI(repo.URL); err != nil {
		return errors.New("URL repository tidak valid")
	}
	if repo.UserID == 0 {
		return errors.New("user ID tidak boleh kosong")
	}
	return nil
}
