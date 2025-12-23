//go:build integration

package repository_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"bookhub/internal/domain/entity"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Exported variables for use in other test files
var (
	MongoTestDB    *mongo.Database
	PostgresTestDB *sql.DB
)

// Shared containers and connections for all tests
var (
	// MongoDB
	mongoContainer *mongodb.MongoDBContainer
	mongoClient    *mongo.Client
	mongoTestDB    *mongo.Database

	// PostgreSQL
	postgresContainer *postgres.PostgresContainer
	postgresDB        *sql.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Setup MongoDB
	if err := setupMongoDB(ctx); err != nil {
		panic("failed to setup MongoDB: " + err.Error())
	}

	// Setup PostgreSQL
	if err := setupPostgres(ctx); err != nil {
		panic("failed to setup PostgreSQL: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup
	teardown(ctx)

	os.Exit(code)
}

func setupMongoDB(ctx context.Context) error {
	var err error
	mongoContainer, err = mongodb.Run(ctx, "mongo:7")
	if err != nil {
		return err
	}

	connectionString, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		return err
	}

	clientOptions := options.Client().ApplyURI(connectionString)
	mongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	if err = mongoClient.Ping(ctx, nil); err != nil {
		return err
	}

	mongoTestDB = mongoClient.Database("bookhub_test")
	MongoTestDB = mongoTestDB
	return nil
}

func setupPostgres(ctx context.Context) error {
	var err error
	postgresContainer, err = postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("bookhub_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return err
	}

	connectionString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return err
	}

	postgresDB, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	// Retry ping with backoff
	for i := 0; i < 10; i++ {
		if err = postgresDB.Ping(); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return err
	}

	// Run migrations
	if err = runPostgresMigrations(postgresDB); err != nil {
		return err
	}

	PostgresTestDB = postgresDB
	return nil
}

func runPostgresMigrations(db *sql.DB) error {
	migrations := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(100) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,

		// Books table
		`CREATE TABLE IF NOT EXISTS books (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title VARCHAR(200) NOT NULL,
			author VARCHAR(100) NOT NULL,
			isbn VARCHAR(13) NOT NULL UNIQUE,
			published_year INTEGER,
			total_copies INTEGER NOT NULL DEFAULT 1,
			available_copies INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			CONSTRAINT chk_copies CHECK (available_copies >= 0 AND available_copies <= total_copies)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn)`,

		// Loans table
		`CREATE TABLE IF NOT EXISTS loans (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
			borrowed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			due_date TIMESTAMP WITH TIME ZONE NOT NULL,
			returned_at TIMESTAMP WITH TIME ZONE,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			CONSTRAINT chk_status CHECK (status IN ('active', 'returned'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_loans_user_id ON loans(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_loans_book_id ON loans(book_id)`,
		`CREATE INDEX IF NOT EXISTS idx_loans_status ON loans(status)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return err
		}
	}

	return nil
}

func teardown(ctx context.Context) {
	if mongoClient != nil {
		_ = mongoClient.Disconnect(ctx)
	}
	if mongoContainer != nil {
		_ = mongoContainer.Terminate(ctx)
	}
	if postgresDB != nil {
		_ = postgresDB.Close()
	}
	if postgresContainer != nil {
		_ = postgresContainer.Terminate(ctx)
	}
}

// CleanupMongo clears all MongoDB collections between tests
func CleanupMongo(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	_ = mongoTestDB.Collection("books").Drop(ctx)
	_ = mongoTestDB.Collection("users").Drop(ctx)
	_ = mongoTestDB.Collection("loans").Drop(ctx)
}

// CleanupPostgres clears all PostgreSQL tables between tests
func CleanupPostgres(t *testing.T) {
	t.Helper()
	// Delete in correct order due to foreign key constraints
	_, _ = postgresDB.Exec("DELETE FROM loans")
	_, _ = postgresDB.Exec("DELETE FROM books")
	_, _ = postgresDB.Exec("DELETE FROM users")
}

// ============================================================================
// Test Helper Functions
// ============================================================================

func CreateTestBook(title, author, isbn string) *entity.Book {
	return &entity.Book{
		ID:              uuid.New(),
		Title:           title,
		Author:          author,
		ISBN:            isbn,
		PublishedYear:   2024,
		TotalCopies:     5,
		AvailableCopies: 5,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func CreateTestUser(name, email string) *entity.User {
	return &entity.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: "hashedpassword123",
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func CreateTestLoan(userID, bookID uuid.UUID) *entity.Loan {
	return &entity.Loan{
		ID:         uuid.New(),
		UserID:     userID,
		BookID:     bookID,
		BorrowedAt: time.Now(),
		DueDate:    time.Now().AddDate(0, 0, 14),
		ReturnedAt: nil,
		Status:     entity.LoanStatusActive,
	}
}
