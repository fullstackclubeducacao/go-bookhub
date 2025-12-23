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

func TestPostgresBookRepository_Create(t *testing.T) {
	CleanupPostgres(t)

	repo := repository.NewPostgresBookRepository(PostgresTestDB)
	ctx := context.Background()

	book := CreateTestBook("Test Book PG", "Test Author", "1234567890")

	err := repo.Create(ctx, book)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, book.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, book.Title, retrieved.Title)
	assert.Equal(t, book.Author, retrieved.Author)
	assert.Equal(t, book.ISBN, retrieved.ISBN)
}

func TestPostgresBookRepository_GetByID(t *testing.T) {
	CleanupPostgres(t)

	repo := repository.NewPostgresBookRepository(PostgresTestDB)
	ctx := context.Background()

	book := CreateTestBook("GetByID Book PG", "Author", "1234567891")
	err := repo.Create(ctx, book)
	require.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, book.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, book.ID, retrieved.ID)

	nonExistent, err := repo.GetByID(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, nonExistent)
}

func TestPostgresBookRepository_GetByISBN(t *testing.T) {
	CleanupPostgres(t)

	repo := repository.NewPostgresBookRepository(PostgresTestDB)
	ctx := context.Background()

	book := CreateTestBook("GetByISBN Book PG", "Author", "9876543210")
	err := repo.Create(ctx, book)
	require.NoError(t, err)

	retrieved, err := repo.GetByISBN(ctx, book.ISBN)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, book.ISBN, retrieved.ISBN)

	nonExistent, err := repo.GetByISBN(ctx, "0000000000")
	assert.NoError(t, err)
	assert.Nil(t, nonExistent)
}

func TestPostgresBookRepository_List(t *testing.T) {
	CleanupPostgres(t)

	repo := repository.NewPostgresBookRepository(PostgresTestDB)
	ctx := context.Background()

	books := []*entity.Book{
		CreateTestBook("Book 1 PG", "Author 1", "1111111111"),
		CreateTestBook("Book 2 PG", "Author 2", "2222222222"),
		CreateTestBook("Book 3 PG", "Author 3", "3333333333"),
	}
	books[2].AvailableCopies = 0

	for _, book := range books {
		err := repo.Create(ctx, book)
		require.NoError(t, err)
	}

	result, count, err := repo.List(ctx, 1, 10, nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
	assert.Len(t, result, 3)

	result, count, err = repo.List(ctx, 1, 2, nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
	assert.Len(t, result, 2)

	availableOnly := true
	result, count, err = repo.List(ctx, 1, 10, &availableOnly)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, result, 2)
}

func TestPostgresBookRepository_Update(t *testing.T) {
	CleanupPostgres(t)

	repo := repository.NewPostgresBookRepository(PostgresTestDB)
	ctx := context.Background()

	book := CreateTestBook("Original Title PG", "Original Author", "1234567892")
	err := repo.Create(ctx, book)
	require.NoError(t, err)

	book.Title = "Updated Title PG"
	book.Author = "Updated Author"
	book.AvailableCopies = 3
	book.UpdatedAt = time.Now()

	err = repo.Update(ctx, book)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, book.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title PG", retrieved.Title)
	assert.Equal(t, "Updated Author", retrieved.Author)
	assert.Equal(t, 3, retrieved.AvailableCopies)
}

func TestPostgresBookRepository_Delete(t *testing.T) {
	CleanupPostgres(t)

	repo := repository.NewPostgresBookRepository(PostgresTestDB)
	ctx := context.Background()

	book := CreateTestBook("Delete Book PG", "Author", "1234567893")
	err := repo.Create(ctx, book)
	require.NoError(t, err)

	err = repo.Delete(ctx, book.ID)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, book.ID)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}
