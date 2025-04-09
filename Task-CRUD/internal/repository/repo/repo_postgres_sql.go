package repo

import (
	"Task-CRUD/internal/entity"
	"context"
	"database/sql"
)

// RepoRepositoryInterface mendefinisikan kontrak fungsi untuk Repository
type RepoRepositoryInterface interface {
	GetAllRepositories() ([]entity.Repository, error)
	GetRepositoryByID(id uint) (*entity.Repository, error)
	CreateRepository(repo *entity.Repository) error
	UpdateRepository(id uint, updatedRepo *entity.Repository) error
	DeleteRepository(id uint) error
}

// Implementasi PostgreSQL
type RepoRepositoryPostgres struct {
	db *sql.DB
}

func NewRepoRepositoryPostgres(db *sql.DB) RepoRepositoryInterface {
	return &RepoRepositoryPostgres{db: db}
}

func (r *RepoRepositoryPostgres) GetAllRepositories() ([]entity.Repository, error) {
	query := `
	SELECT r.id, r.name, r.user_id, r.url, r.ai_enabled, r.created_at, r.updated_at,
	       u.id, u.name, u.email, u.created_at, u.updated_at
	FROM repositories r
	JOIN users u ON r.user_id = u.id
	`
	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []entity.Repository
	for rows.Next() {
		var repo entity.Repository
		var user entity.User
		err := rows.Scan(
			&repo.ID, &repo.Name, &repo.UserID, &repo.URL, &repo.AIEnabled, &repo.CreatedAt, &repo.UpdatedAt,
			&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		repo.User = user
		repos = append(repos, repo)
	}

	return repos, nil
}

func (r *RepoRepositoryPostgres) GetRepositoryByID(id uint) (*entity.Repository, error) {
	query := `
	SELECT r.id, r.name, r.user_id, r.url, r.ai_enabled, r.created_at, r.updated_at,
	       u.id, u.name, u.email, u.created_at, u.updated_at
	FROM repositories r
	JOIN users u ON r.user_id = u.id
	WHERE r.id = $1
	`
	var repo entity.Repository
	var user entity.User

	err := r.db.QueryRowContext(context.Background(), query, id).Scan(
		&repo.ID, &repo.Name, &repo.UserID, &repo.URL, &repo.AIEnabled, &repo.CreatedAt, &repo.UpdatedAt,
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	repo.User = user
	return &repo, nil
}

func (r *RepoRepositoryPostgres) CreateRepository(repo *entity.Repository) error {
	query := `
	INSERT INTO repositories (name, user_id, url, ai_enabled, created_at, updated_at)
	VALUES ($1, $2, $3, $4, NOW(), NOW())
	RETURNING id
	`
	return r.db.QueryRowContext(context.Background(), query,
		repo.Name, repo.UserID, repo.URL, repo.AIEnabled,
	).Scan(&repo.ID)
}

func (r *RepoRepositoryPostgres) UpdateRepository(id uint, updatedRepo *entity.Repository) error {
	query := `
	UPDATE repositories
	SET name = $1, user_id = $2, url = $3, ai_enabled = $4, updated_at = NOW()
	WHERE id = $5
	`
	_, err := r.db.ExecContext(context.Background(), query,
		updatedRepo.Name, updatedRepo.UserID, updatedRepo.URL, updatedRepo.AIEnabled, id,
	)
	return err
}

func (r *RepoRepositoryPostgres) DeleteRepository(id uint) error {
	_, err := r.db.ExecContext(context.Background(), `DELETE FROM repositories WHERE id = $1`, id)
	return err
}
