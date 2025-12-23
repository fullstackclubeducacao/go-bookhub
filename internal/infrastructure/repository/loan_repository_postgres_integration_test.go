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

func TestPostgresLoanRepository_Create(t *testing.T) {
	CleanupPostgres(t)

	ctx := context.Background()
	userRepo := repository.NewPostgresUserRepository(PostgresTestDB)
	bookRepo := repository.NewPostgresBookRepository(PostgresTestDB)
	repo := repository.NewPostgresLoanRepository(PostgresTestDB)

	user := CreateTestUser("Loan User PG", "loanuserpg@example.com")
	book := CreateTestBook("Loan Book PG", "Author", "1234567894")

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	err = bookRepo.Create(ctx, book)
	require.NoError(t, err)

	loan := CreateTestLoan(user.ID, book.ID)
	err = repo.Create(ctx, loan)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, loan.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, loan.UserID, retrieved.UserID)
	assert.Equal(t, loan.BookID, retrieved.BookID)
	assert.Equal(t, entity.LoanStatusActive, retrieved.Status)
}

func TestPostgresLoanRepository_GetByID(t *testing.T) {
	CleanupPostgres(t)

	ctx := context.Background()
	userRepo := repository.NewPostgresUserRepository(PostgresTestDB)
	bookRepo := repository.NewPostgresBookRepository(PostgresTestDB)
	repo := repository.NewPostgresLoanRepository(PostgresTestDB)

	user := CreateTestUser("GetByID Loan User PG", "getbyidloanpg@example.com")
	book := CreateTestBook("GetByID Loan Book PG", "Author", "1234567895")

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	err = bookRepo.Create(ctx, book)
	require.NoError(t, err)

	loan := CreateTestLoan(user.ID, book.ID)
	err = repo.Create(ctx, loan)
	require.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, loan.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, loan.ID, retrieved.ID)

	nonExistent, err := repo.GetByID(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, nonExistent)
}

func TestPostgresLoanRepository_GetByIDWithDetails(t *testing.T) {
	CleanupPostgres(t)

	ctx := context.Background()
	userRepo := repository.NewPostgresUserRepository(PostgresTestDB)
	bookRepo := repository.NewPostgresBookRepository(PostgresTestDB)
	repo := repository.NewPostgresLoanRepository(PostgresTestDB)

	user := CreateTestUser("Details User PG", "detailsuserpg@example.com")
	book := CreateTestBook("Details Book PG", "Details Author", "1234567896")

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	err = bookRepo.Create(ctx, book)
	require.NoError(t, err)

	loan := CreateTestLoan(user.ID, book.ID)
	err = repo.Create(ctx, loan)
	require.NoError(t, err)

	details, err := repo.GetByIDWithDetails(ctx, loan.ID)
	assert.NoError(t, err)
	assert.NotNil(t, details)
	assert.Equal(t, loan.ID, details.Loan.ID)
	assert.Equal(t, user.Name, details.UserName)
	assert.Equal(t, book.Title, details.BookTitle)

	nonExistent, err := repo.GetByIDWithDetails(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, nonExistent)
}

func TestPostgresLoanRepository_GetActiveByUserAndBook(t *testing.T) {
	CleanupPostgres(t)

	ctx := context.Background()
	userRepo := repository.NewPostgresUserRepository(PostgresTestDB)
	bookRepo := repository.NewPostgresBookRepository(PostgresTestDB)
	repo := repository.NewPostgresLoanRepository(PostgresTestDB)

	user := CreateTestUser("Active Loan User PG", "activeloanpg@example.com")
	book := CreateTestBook("Active Loan Book PG", "Author", "1234567897")

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	err = bookRepo.Create(ctx, book)
	require.NoError(t, err)

	loan := CreateTestLoan(user.ID, book.ID)
	err = repo.Create(ctx, loan)
	require.NoError(t, err)

	active, err := repo.GetActiveByUserAndBook(ctx, user.ID, book.ID)
	assert.NoError(t, err)
	assert.NotNil(t, active)
	assert.Equal(t, loan.ID, active.ID)

	nonExistent, err := repo.GetActiveByUserAndBook(ctx, uuid.New(), uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, nonExistent)
}

func TestPostgresLoanRepository_List(t *testing.T) {
	CleanupPostgres(t)

	ctx := context.Background()
	userRepo := repository.NewPostgresUserRepository(PostgresTestDB)
	bookRepo := repository.NewPostgresBookRepository(PostgresTestDB)
	repo := repository.NewPostgresLoanRepository(PostgresTestDB)

	user1 := CreateTestUser("List User 1 PG", "listuser1pg@example.com")
	user2 := CreateTestUser("List User 2 PG", "listuser2pg@example.com")
	book1 := CreateTestBook("List Book 1 PG", "Author", "1111111112")
	book2 := CreateTestBook("List Book 2 PG", "Author", "2222222223")
	book3 := CreateTestBook("List Book 3 PG", "Author", "3333333334")

	require.NoError(t, userRepo.Create(ctx, user1))
	require.NoError(t, userRepo.Create(ctx, user2))
	require.NoError(t, bookRepo.Create(ctx, book1))
	require.NoError(t, bookRepo.Create(ctx, book2))
	require.NoError(t, bookRepo.Create(ctx, book3))

	loan1 := CreateTestLoan(user1.ID, book1.ID)
	loan2 := CreateTestLoan(user1.ID, book2.ID)
	loan3 := CreateTestLoan(user2.ID, book3.ID)
	loan3.Status = entity.LoanStatusReturned
	now := time.Now()
	loan3.ReturnedAt = &now

	require.NoError(t, repo.Create(ctx, loan1))
	require.NoError(t, repo.Create(ctx, loan2))
	require.NoError(t, repo.Create(ctx, loan3))

	result, count, err := repo.List(ctx, 1, 10, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
	assert.Len(t, result, 3)

	result, count, err = repo.List(ctx, 1, 10, &user1.ID, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, result, 2)

	activeStatus := entity.LoanStatusActive
	result, count, err = repo.List(ctx, 1, 10, nil, &activeStatus)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	returnedStatus := entity.LoanStatusReturned
	result, count, err = repo.List(ctx, 1, 10, nil, &returnedStatus)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestPostgresLoanRepository_ListWithDetails(t *testing.T) {
	CleanupPostgres(t)

	ctx := context.Background()
	userRepo := repository.NewPostgresUserRepository(PostgresTestDB)
	bookRepo := repository.NewPostgresBookRepository(PostgresTestDB)
	repo := repository.NewPostgresLoanRepository(PostgresTestDB)

	user := CreateTestUser("ListDetails User PG", "listdetailspg@example.com")
	book1 := CreateTestBook("ListDetails Book 1 PG", "Author 1", "4444444444")
	book2 := CreateTestBook("ListDetails Book 2 PG", "Author 2", "5555555555")

	require.NoError(t, userRepo.Create(ctx, user))
	require.NoError(t, bookRepo.Create(ctx, book1))
	require.NoError(t, bookRepo.Create(ctx, book2))

	loan1 := CreateTestLoan(user.ID, book1.ID)
	loan2 := CreateTestLoan(user.ID, book2.ID)

	require.NoError(t, repo.Create(ctx, loan1))
	require.NoError(t, repo.Create(ctx, loan2))

	details, count, err := repo.ListWithDetails(ctx, 1, 10, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, details, 2)

	for _, d := range details {
		assert.Equal(t, user.Name, d.UserName)
		assert.NotEmpty(t, d.BookTitle)
	}
}

func TestPostgresLoanRepository_Update(t *testing.T) {
	CleanupPostgres(t)

	ctx := context.Background()
	userRepo := repository.NewPostgresUserRepository(PostgresTestDB)
	bookRepo := repository.NewPostgresBookRepository(PostgresTestDB)
	repo := repository.NewPostgresLoanRepository(PostgresTestDB)

	user := CreateTestUser("Update Loan User PG", "updateloanpg@example.com")
	book := CreateTestBook("Update Loan Book PG", "Author", "1234567898")

	require.NoError(t, userRepo.Create(ctx, user))
	require.NoError(t, bookRepo.Create(ctx, book))

	loan := CreateTestLoan(user.ID, book.ID)
	require.NoError(t, repo.Create(ctx, loan))

	now := time.Now()
	loan.ReturnedAt = &now
	loan.Status = entity.LoanStatusReturned

	err := repo.Update(ctx, loan)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, loan.ID)
	assert.NoError(t, err)
	assert.Equal(t, entity.LoanStatusReturned, retrieved.Status)
	assert.NotNil(t, retrieved.ReturnedAt)
}
