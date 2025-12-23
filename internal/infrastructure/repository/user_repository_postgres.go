package repository

import (
	"context"
	"database/sql"

	"bookhub/internal/domain/entity"
	"bookhub/internal/domain/repository"
	"bookhub/internal/infrastructure/database/sqlc"

	"github.com/google/uuid"
)

type postgresUserRepository struct {
	queries *sqlc.Queries
}

func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &postgresUserRepository{
		queries: sqlc.New(db),
	}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	_, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Active:       user.Active,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	})
	return err
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(row), nil
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(row), nil
}

func (r *postgresUserRepository) List(ctx context.Context, page, limit int) ([]*entity.User, int, error) {
	offset := (page - 1) * limit

	rows, err := r.queries.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	count, err := r.queries.CountUsers(ctx)
	if err != nil {
		return nil, 0, err
	}

	users := make([]*entity.User, len(rows))
	for i, row := range rows {
		users[i] = r.toEntity(row)
	}

	return users, int(count), nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	_, err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Active:    user.Active,
		UpdatedAt: user.UpdatedAt,
	})
	return err
}

func (r *postgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *postgresUserRepository) toEntity(row sqlc.User) *entity.User {
	return &entity.User{
		ID:           row.ID,
		Name:         row.Name,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Active:       row.Active,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
