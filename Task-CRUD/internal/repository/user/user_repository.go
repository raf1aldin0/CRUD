package user

import (
	"Task-CRUD/internal/entity"
	"log"
	"time"

	"gorm.io/gorm"
)

// Interface langsung didefinisikan di file ini
type UserRepositoryInterfaceGorm interface {
	GetAllUsers() ([]entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(id uint, user *entity.User) error
	DeleteUser(id uint) error
}

type UserRepositoryGorm struct {
	db *gorm.DB
}

func NewUserRepositoryGorm(db *gorm.DB) UserRepositoryInterfaceGorm {
	return &UserRepositoryGorm{db: db}
}

func (r *UserRepositoryGorm) GetAllUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepositoryGorm) GetUserByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepositoryGorm) CreateUser(user *entity.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := r.db.Create(user).Error
	if err != nil {
		log.Printf("ERROR | GORM gagal insert user: %v", err)
	}
	return err
}

func (r *UserRepositoryGorm) UpdateUser(id uint, user *entity.User) error {
	user.UpdatedAt = time.Now()
	return r.db.Model(&entity.User{}).Where("id = ?", id).Updates(user).Error
}

func (r *UserRepositoryGorm) DeleteUser(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}
