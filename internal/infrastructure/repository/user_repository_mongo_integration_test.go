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

func TestMongoUserRepository_Create(t *testing.T) {
	CleanupMongo(t)

	repo := repository.NewMongoUserRepository(MongoTestDB)
	ctx := context.Background()

	user := CreateTestUser("Test User", "test@example.com")

	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, user.Name, retrieved.Name)
	assert.Equal(t, user.Email, retrieved.Email)
}

func TestMongoUserRepository_GetByID(t *testing.T) {
	CleanupMongo(t)

	repo := repository.NewMongoUserRepository(MongoTestDB)
	ctx := context.Background()

	user := CreateTestUser("GetByID User", "getbyid@example.com")
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

func TestMongoUserRepository_GetByEmail(t *testing.T) {
	CleanupMongo(t)

	repo := repository.NewMongoUserRepository(MongoTestDB)
	ctx := context.Background()

	user := CreateTestUser("GetByEmail User", "getbyemail@example.com")
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	retrieved, err := repo.GetByEmail(ctx, user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, user.Email, retrieved.Email)

	nonExistent, err := repo.GetByEmail(ctx, "nonexistent@example.com")
	assert.NoError(t, err)
	assert.Nil(t, nonExistent)
}

func TestMongoUserRepository_List(t *testing.T) {
	CleanupMongo(t)

	repo := repository.NewMongoUserRepository(MongoTestDB)
	ctx := context.Background()

	users := []*entity.User{
		CreateTestUser("User 1", "user1@example.com"),
		CreateTestUser("User 2", "user2@example.com"),
		CreateTestUser("User 3", "user3@example.com"),
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

func TestMongoUserRepository_Update(t *testing.T) {
	CleanupMongo(t)

	repo := repository.NewMongoUserRepository(MongoTestDB)
	ctx := context.Background()

	user := CreateTestUser("Original Name", "original@example.com")
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	user.Name = "Updated Name"
	user.Email = "updated@example.com"
	user.Active = false
	user.UpdatedAt = time.Now()

	err = repo.Update(ctx, user)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", retrieved.Name)
	assert.Equal(t, "updated@example.com", retrieved.Email)
	assert.False(t, retrieved.Active)
}

func TestMongoUserRepository_Delete(t *testing.T) {
	CleanupMongo(t)

	repo := repository.NewMongoUserRepository(MongoTestDB)
	ctx := context.Background()

	user := CreateTestUser("Delete User", "delete@example.com")
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	err = repo.Delete(ctx, user.ID)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}
