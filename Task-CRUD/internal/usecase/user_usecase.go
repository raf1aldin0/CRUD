package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"Task-CRUD/internal/entity"
	"Task-CRUD/internal/repository/user" // ⬅️ ganti import ke package user

	"github.com/redis/go-redis/v9"
)

type UserUseCaseInterface interface {
	GetUsers(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, id uint, user *entity.User) error
	DeleteUser(ctx context.Context, id uint) error
}

type UserUseCase struct {
	userRepo user.UserRepositoryInterfaceGorm // ⬅️ pakai interface dari package user
	redis    *redis.Client
}

// Constructor tanpa Redis
func NewUserUseCase(userRepo user.UserRepositoryInterfaceGorm) UserUseCaseInterface {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// Constructor dengan Redis
func NewUserUseCaseWithCache(userRepo user.UserRepositoryInterfaceGorm, redisClient *redis.Client) UserUseCaseInterface {
	return &UserUseCase{
		userRepo: userRepo,
		redis:    redisClient,
	}
}

// GetUsers mengambil semua user dan cache ke Redis jika ada
func (uc *UserUseCase) GetUsers(ctx context.Context) ([]entity.User, error) {
	cacheKey := "users:all"

	// Ambil dari Redis dulu
	if uc.redis != nil {
		cached, err := uc.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var users []entity.User
			if err := json.Unmarshal([]byte(cached), &users); err == nil {
				fmt.Println("✅ Data users diambil dari Redis")
				return users, nil
			}
			fmt.Printf("⚠️ Gagal unmarshal data users dari Redis: %v\n", err)
		} else if err != redis.Nil {
			fmt.Printf("⚠️ Redis error: %v\n", err)
		}
	}

	// Ambil dari database
	users, err := uc.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	// Simpan ke Redis
	if uc.redis != nil {
		data, _ := json.Marshal(users)
		if err := uc.redis.Set(ctx, cacheKey, data, 10*time.Minute).Err(); err != nil {
			fmt.Printf("⚠️ Gagal set Redis users: %v\n", err)
		}
	}

	return users, nil
}

// GetUserByID mengambil user berdasarkan ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	return uc.userRepo.GetUserByID(id)
}

// CreateUser membuat user baru dan invalidate cache
func (uc *UserUseCase) CreateUser(ctx context.Context, user *entity.User) error {
	if user.Name == "" {
		return errors.New("nama tidak boleh kosong")
	}
	if user.Email == "" {
		return errors.New("email tidak boleh kosong")
	}

	if err := uc.userRepo.CreateUser(user); err != nil {
		return err
	}

	if uc.redis != nil {
		if err := uc.redis.Del(ctx, "users:all").Err(); err != nil {
			fmt.Printf("⚠️ Gagal hapus cache users setelah Create: %v\n", err)
		}
	}

	return nil
}

// UpdateUser memperbarui data user dan invalidate cache
func (uc *UserUseCase) UpdateUser(ctx context.Context, id uint, user *entity.User) error {
	if user.Name == "" {
		return errors.New("nama tidak boleh kosong")
	}

	if err := uc.userRepo.UpdateUser(id, user); err != nil {
		return err
	}

	if uc.redis != nil {
		if err := uc.redis.Del(ctx, "users:all").Err(); err != nil {
			fmt.Printf("⚠️ Gagal hapus cache users setelah Update: %v\n", err)
		}
	}

	return nil
}

// DeleteUser menghapus user dan invalidate cache
func (uc *UserUseCase) DeleteUser(ctx context.Context, id uint) error {
	if err := uc.userRepo.DeleteUser(id); err != nil {
		return err
	}

	if uc.redis != nil {
		if err := uc.redis.Del(ctx, "users:all").Err(); err != nil {
			fmt.Printf("⚠️ Gagal hapus cache users setelah Delete: %v\n", err)
		}
	}

	return nil
}
