package user

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

type UserRepositoryGorm struct {
	db *gorm.DB
}

func NewUserRepositoryGorm(db *gorm.DB) interfaces.UserRepositoryInterfaceGorm {
	return &UserRepositoryGorm{db: db}
}

func (r *UserRepositoryGorm) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepositoryGorm.GetAllUsers")
	defer span.Finish()

	var users []entity.User
	err := r.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		ext.LogError(span, err)
	}
	return users, err
}

func (r *UserRepositoryGorm) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepositoryGorm.GetUserByID")
	defer span.Finish()

	var user entity.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		ext.LogError(span, err)
	}
	return &user, err
}

func (r *UserRepositoryGorm) CreateUser(ctx context.Context, user *entity.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepositoryGorm.CreateUser")
	defer span.Finish()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		ext.LogError(span, err)
		log.Printf("ERROR | GORM gagal insert user: %v", err)
	}
	return err
}

func (r *UserRepositoryGorm) UpdateUser(ctx context.Context, id uint, user *entity.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepositoryGorm.UpdateUser")
	defer span.Finish()

	user.UpdatedAt = time.Now()
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Updates(user).Error
	if err != nil {
		ext.LogError(span, err)
	}
	return err
}

func (r *UserRepositoryGorm) DeleteUser(ctx context.Context, id uint) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepositoryGorm.DeleteUser")
	defer span.Finish()

	err := r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
	if err != nil {
		ext.LogError(span, err)
	}
	return err
}
