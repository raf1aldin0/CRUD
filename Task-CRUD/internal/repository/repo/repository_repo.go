package repo

import (
	"Task-CRUD/internal/entity"
	interfaces "Task-CRUD/internal/interfaces"
	"context"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gorm.io/gorm"
)

type RepoRepositoryGorm struct {
	db *gorm.DB
}

func NewRepoRepositoryGorm(db *gorm.DB) interfaces.RepoRepositoryInterfaceGorm {
	return &RepoRepositoryGorm{db: db}
}

func (r *RepoRepositoryGorm) GetAllRepositories(ctx context.Context) ([]entity.Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoRepository.GetAllRepositories")
	defer span.Finish()

	var repos []entity.Repository
	if err := r.db.WithContext(ctx).Preload("User").Find(&repos).Error; err != nil {
		ext.LogError(span, err)
		return nil, err
	}
	return repos, nil
}

func (r *RepoRepositoryGorm) GetRepositoryByID(ctx context.Context, id uint) (*entity.Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoRepository.GetRepositoryByID")
	defer span.Finish()

	var repo entity.Repository
	if err := r.db.WithContext(ctx).Preload("User").First(&repo, id).Error; err != nil {
		ext.LogError(span, err)
		return nil, err
	}
	return &repo, nil
}

func (r *RepoRepositoryGorm) CreateRepository(ctx context.Context, repo *entity.Repository) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoRepository.CreateRepository")
	defer span.Finish()

	repo.CreatedAt = time.Now()
	repo.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Create(repo).Error; err != nil {
		log.Printf("ERROR | GORM gagal insert repository: %v", err)
		ext.LogError(span, err)
		return err
	}
	return nil
}

func (r *RepoRepositoryGorm) UpdateRepository(ctx context.Context, id uint, updatedRepo *entity.Repository) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoRepository.UpdateRepository")
	defer span.Finish()

	updatedRepo.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Model(&entity.Repository{}).Where("id = ?", id).Updates(updatedRepo).Error; err != nil {
		ext.LogError(span, err)
		return err
	}
	return nil
}

func (r *RepoRepositoryGorm) DeleteRepository(ctx context.Context, id uint) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoRepository.DeleteRepository")
	defer span.Finish()

	if err := r.db.WithContext(ctx).Delete(&entity.Repository{}, id).Error; err != nil {
		ext.LogError(span, err)
		return err
	}
	return nil
}
