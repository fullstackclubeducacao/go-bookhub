//go:build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"bookhub/internal/domain/entity"
	"bookhub/internal/infrastructure/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresUserRepository_Create(t *testing.T) {
	CleanupPostgres(t)

	repo := repository.NewPostgresUserRepository(PostgresTestDB)
	ctx := context.Background()

	user := CreateTestUser("Test User PG", "testpg@example.com")

	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, user.Name, retrieved.Name)
	assert.Equal(t, user.Email, retrieved.Email)
}

func TestPostgresUserRepository_GetByID(t *testing.T) {
	CleanupPostgres(t)

	repo := repository.NewPostgresUserRepository(PostgresTestDB)
	ctx := context.Background()

	user := CreateTestUser("GetByID User PG", "getbyidpg@example.com")
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, user.ID, retrieved.ID)

	nonExistent, err := repo.GetByID(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, nonExistent)
}

func TestPostgresUserRepository_GetByEmail(t *testing.T) {
	CleanupPostgres(t)

	repo := repository.NewPostgresUserRepository(PostgresTestDB)
	ctx := context.Background()

	user := CreateTestUser("GetByEmail User PG", "getbyemailpg@example.com")
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	retrieved, err := repo.GetByEmail(ctx, user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, user.Email, retrieved.Email)

	nonExistent, err := repo.GetByEmail(ctx, "nonexistentpg@example.com")
	assert.NoError(t, err)
	assert.Nil(t, nonExistent)
}

func TestPostgresUserRepository_List(t *testing.T) {
	CleanupPostgres(t)

	repo := repository.NewPostgresUserRepository(PostgresTestDB)
	ctx := context.Background()

	users := []*entity.User{
		CreateTestUser("User 1 PG", "user1pg@example.com"),
		CreateTestUser("User 2 PG", "user2pg@example.com"),
		CreateTestUser("User 3 PG", "user3pg@example.com"),
	}

	for _, user := range users {
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	result, count, err := repo.List(ctx, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
	assert.Len(t, result, 3)

	result, count, err = repo.List(ctx, 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
	assert.Len(t, result, 2)
}

func TestPostgresUserRepository_Update(t *testing.T) {
	CleanupPostgres(t)

	repo := repository.NewPostgresUserRepository(PostgresTestDB)
	ctx := context.Background()

	user := CreateTestUser("Original Name PG", "originalpg@example.com")
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	user.Name = "Updated Name PG"
	user.Email = "updatedpg@example.com"
	user.Active = false
	user.UpdatedAt = time.Now()

	err = repo.Update(ctx, user)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name PG", retrieved.Name)
	assert.Equal(t, "updatedpg@example.com", retrieved.Email)
	assert.False(t, retrieved.Active)
}

func TestPostgresUserRepository_Delete(t *testing.T) {
	CleanupPostgres(t)

	repo := repository.NewPostgresUserRepository(PostgresTestDB)
	ctx := context.Background()

	user := CreateTestUser("Delete User PG", "deletepg@example.com")
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	err = repo.Delete(ctx, user.ID)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}
