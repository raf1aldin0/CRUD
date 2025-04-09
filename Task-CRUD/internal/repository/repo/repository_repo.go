package repo

import (
	"Task-CRUD/internal/entity"
	"log"
	"time"

	"gorm.io/gorm"
)

// RepoRepositoryInterfaceGorm mendefinisikan kontrak fungsi Repository
type RepoRepositoryInterfaceGorm interface {
	GetAllRepositories() ([]entity.Repository, error)
	GetRepositoryByID(id uint) (*entity.Repository, error)
	CreateRepository(repo *entity.Repository) error
	UpdateRepository(id uint, updatedRepo *entity.Repository) error
	DeleteRepository(id uint) error
}

type RepoRepositoryGorm struct {
	db *gorm.DB
}

func NewRepoRepositoryGorm(db *gorm.DB) RepoRepositoryInterfaceGorm {
	return &RepoRepositoryGorm{db: db}
}

func (r *RepoRepositoryGorm) GetAllRepositories() ([]entity.Repository, error) {
	var repos []entity.Repository
	if err := r.db.Preload("User").Find(&repos).Error; err != nil {
		return nil, err
	}
	return repos, nil
}

func (r *RepoRepositoryGorm) GetRepositoryByID(id uint) (*entity.Repository, error) {
	var repo entity.Repository
	if err := r.db.Preload("User").First(&repo, id).Error; err != nil {
		return nil, err
	}
	return &repo, nil
}

func (r *RepoRepositoryGorm) CreateRepository(repo *entity.Repository) error {
	repo.CreatedAt = time.Now()
	repo.UpdatedAt = time.Now()
	if err := r.db.Create(repo).Error; err != nil {
		log.Printf("ERROR | GORM gagal insert repository: %v", err)
		return err
	}
	return nil
}

func (r *RepoRepositoryGorm) UpdateRepository(id uint, updatedRepo *entity.Repository) error {
	updatedRepo.UpdatedAt = time.Now()
	return r.db.Model(&entity.Repository{}).Where("id = ?", id).Updates(updatedRepo).Error
}

func (r *RepoRepositoryGorm) DeleteRepository(id uint) error {
	return r.db.Delete(&entity.Repository{}, id).Error
}
