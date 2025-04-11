package user

import (
	"Task-CRUD/internal/entity"
	"context"
	"database/sql"
	interfaces "Task-CRUD/internal/interfaces"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(db *sql.DB) interfaces.UserRepositoryInterfaceSQL {
	return &UserRepositoryPostgres{db: db}
}

func (r *UserRepositoryPostgres) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepositoryPostgres.GetAllUsers")
	defer span.Finish()

	rows, err := r.db.QueryContext(ctx, `SELECT id, name, email, created_at, updated_at FROM users`)
	if err != nil {
		ext.LogError(span, err)
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			ext.LogError(span, err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepositoryPostgres) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepositoryPostgres.GetUserByID")
	defer span.Finish()

	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
	var user entity.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		ext.LogError(span, err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryPostgres) CreateUser(ctx context.Context, user *entity.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepositoryPostgres.CreateUser")
	defer span.Finish()

	query := `INSERT INTO users (name, email, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		ext.LogError(span, err)
	}
	return err
}

func (r *UserRepositoryPostgres) UpdateUser(ctx context.Context, id uint, user *entity.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepositoryPostgres.UpdateUser")
	defer span.Finish()

	query := `UPDATE users SET name = $1, email = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, id)
	if err != nil {
		ext.LogError(span, err)
	}
	return err
}

func (r *UserRepositoryPostgres) DeleteUser(ctx context.Context, id uint) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepositoryPostgres.DeleteUser")
	defer span.Finish()

	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		ext.LogError(span, err)
	}
	return err
}
