package user

import (
	"Task-CRUD/internal/entity"
	"context"
	"database/sql"
)

// Interface langsung di-embed di file ini
type UserRepositoryInterfaceSQL interface {
	GetAllUsers() ([]entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(id uint, user *entity.User) error
	DeleteUser(id uint) error
}

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(db *sql.DB) UserRepositoryInterfaceSQL {
	return &UserRepositoryPostgres{db: db}
}

func (r *UserRepositoryPostgres) GetAllUsers() ([]entity.User, error) {
	rows, err := r.db.QueryContext(context.Background(), `SELECT id, name, email, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepositoryPostgres) GetUserByID(id uint) (*entity.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
	var user entity.User
	err := r.db.QueryRowContext(context.Background(), query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryPostgres) CreateUser(user *entity.User) error {
	query := `INSERT INTO users (name, email, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id`
	return r.db.QueryRowContext(context.Background(), query, user.Name, user.Email).Scan(&user.ID)
}

func (r *UserRepositoryPostgres) UpdateUser(id uint, user *entity.User) error {
	query := `UPDATE users SET name = $1, email = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(context.Background(), query, user.Name, user.Email, id)
	return err
}

func (r *UserRepositoryPostgres) DeleteUser(id uint) error {
	_, err := r.db.ExecContext(context.Background(), `DELETE FROM users WHERE id = $1`, id)
	return err
}
