package interfaces

import (
	"context"

"Task-CRUD/internal/entity"
)


// RepoRepositoryInterfaceSQL mendefinisikan kontrak fungsi untuk Repository (SQL)
type RepoRepositoryInterfaceSQL interface {
	GetAllRepositories(ctx context.Context) ([]entity.Repository, error)
	GetRepositoryByID(ctx context.Context, id uint) (*entity.Repository, error)
	CreateRepository(ctx context.Context, repo *entity.Repository) error
	UpdateRepository(ctx context.Context, id uint, updatedRepo *entity.Repository) error
	DeleteRepository(ctx context.Context, id uint) error
}

// RepoRepositoryInterfaceGorm mendefinisikan kontrak fungsi untuk Repository dengan GORM
type RepoRepositoryInterfaceGorm interface {
	GetAllRepositories(ctx context.Context) ([]entity.Repository, error)
	GetRepositoryByID(ctx context.Context, id uint) (*entity.Repository, error)
	CreateRepository(ctx context.Context, repo *entity.Repository) error
	UpdateRepository(ctx context.Context, id uint, updatedRepo *entity.Repository) error
	DeleteRepository(ctx context.Context, id uint) error
}

// UserRepositoryInterfaceSQL mendefinisikan kontrak fungsi untuk User (SQL)
type UserRepositoryInterfaceSQL interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	UpdateUser(ctx context.Context, id uint, user *entity.User) error
	DeleteUser(ctx context.Context, id uint) error
}

// UserRepositoryInterfaceGorm mendefinisikan kontrak fungsi untuk User dengan GORM
type UserRepositoryInterfaceGorm interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	UpdateUser(ctx context.Context, id uint, user *entity.User) error
	DeleteUser(ctx context.Context, id uint) error
}

type RepoUseCaseInterface interface {
	GetAllRepos(ctx context.Context) ([]entity.Repository, error)
	GetRepositoryByID(ctx context.Context, id uint) (*entity.Repository, error)
	CreateRepo(ctx context.Context, repo *entity.Repository) error
	UpdateRepo(ctx context.Context, id uint, repo *entity.Repository) error
	DeleteRepo(ctx context.Context, id uint) error
}

type UserUseCaseInterface interface {
	GetUsers(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, id uint, user *entity.User) error
	DeleteUser(ctx context.Context, id uint) error
}