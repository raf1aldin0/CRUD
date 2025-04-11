package repo

import (
	"context"
	"database/sql"

	"Task-CRUD/internal/entity"
	interfaces "Task-CRUD/internal/interfaces"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type RepoRepositoryPostgres struct {
	db *sql.DB
}

func NewRepoRepositoryPostgres(db *sql.DB) interfaces.RepoRepositoryInterfaceSQL {
	return &RepoRepositoryPostgres{db: db}
}

func (r *RepoRepositoryPostgres) GetAllRepositories(ctx context.Context) ([]entity.Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoRepositorySQL.GetAllRepositories")
	defer span.Finish()

	query := `
	SELECT r.id, r.name, r.user_id, r.url, r.ai_enabled, r.created_at, r.updated_at,
	       u.id, u.name, u.email, u.created_at, u.updated_at
	FROM repositories r
	JOIN users u ON r.user_id = u.id
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		ext.LogError(span, err)
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
			ext.LogError(span, err)
			return nil, err
		}
		repo.User = user
		repos = append(repos, repo)
	}

	return repos, nil
}

func (r *RepoRepositoryPostgres) GetRepositoryByID(ctx context.Context, id uint) (*entity.Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoRepositorySQL.GetRepositoryByID")
	defer span.Finish()

	query := `
	SELECT r.id, r.name, r.user_id, r.url, r.ai_enabled, r.created_at, r.updated_at,
	       u.id, u.name, u.email, u.created_at, u.updated_at
	FROM repositories r
	JOIN users u ON r.user_id = u.id
	WHERE r.id = $1
	`

	var repo entity.Repository
	var user entity.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&repo.ID, &repo.Name, &repo.UserID, &repo.URL, &repo.AIEnabled, &repo.CreatedAt, &repo.UpdatedAt,
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		ext.LogError(span, err)
		return nil, err
	}

	repo.User = user
	return &repo, nil
}

func (r *RepoRepositoryPostgres) CreateRepository(ctx context.Context, repo *entity.Repository) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoRepositorySQL.CreateRepository")
	defer span.Finish()

	query := `
	INSERT INTO repositories (name, user_id, url, ai_enabled, created_at, updated_at)
	VALUES ($1, $2, $3, $4, NOW(), NOW())
	RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query,
		repo.Name, repo.UserID, repo.URL, repo.AIEnabled,
	).Scan(&repo.ID)
	if err != nil {
		ext.LogError(span, err)
	}
	return err
}

func (r *RepoRepositoryPostgres) UpdateRepository(ctx context.Context, id uint, updatedRepo *entity.Repository) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoRepositorySQL.UpdateRepository")
	defer span.Finish()

	query := `
	UPDATE repositories
	SET name = $1, user_id = $2, url = $3, ai_enabled = $4, updated_at = NOW()
	WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, query,
		updatedRepo.Name, updatedRepo.UserID, updatedRepo.URL, updatedRepo.AIEnabled, id,
	)
	if err != nil {
		ext.LogError(span, err)
	}
	return err
}

func (r *RepoRepositoryPostgres) DeleteRepository(ctx context.Context, id uint) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoRepositorySQL.DeleteRepository")
	defer span.Finish()

	_, err := r.db.ExecContext(ctx, `DELETE FROM repositories WHERE id = $1`, id)
	if err != nil {
		ext.LogError(span, err)
	}
	return err
}
