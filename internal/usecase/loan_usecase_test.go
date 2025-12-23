package usecase

import (
	"context"
	"testing"
	"time"

	"bookhub/internal/domain/entity"
	"bookhub/internal/domain/repository"

	"github.com/google/uuid"
)

type mockLoanRepository struct {
	loans map[uuid.UUID]*entity.Loan
}

func newMockLoanRepository() *mockLoanRepository {
	return &mockLoanRepository{
		loans: make(map[uuid.UUID]*entity.Loan),
	}
}

func (m *mockLoanRepository) Create(ctx context.Context, loan *entity.Loan) error {
	m.loans[loan.ID] = loan
	return nil
}

func (m *mockLoanRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Loan, error) {
	if loan, exists := m.loans[id]; exists {
		return loan, nil
	}
	return nil, nil
}

func (m *mockLoanRepository) GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*repository.LoanWithDetails, error) {
	if loan, exists := m.loans[id]; exists {
		return &repository.LoanWithDetails{
			Loan:      loan,
			UserName:  "Test User",
			BookTitle: "Test Book",
		}, nil
	}
	return nil, nil
}

func (m *mockLoanRepository) GetActiveByUserAndBook(ctx context.Context, userID, bookID uuid.UUID) (*entity.Loan, error) {
	for _, loan := range m.loans {
		if loan.UserID == userID && loan.BookID == bookID && loan.Status == entity.LoanStatusActive {
			return loan, nil
		}
	}
	return nil, nil
}

func (m *mockLoanRepository) List(ctx context.Context, page, limit int, userID *uuid.UUID, status *string) ([]*entity.Loan, int, error) {
	loans := make([]*entity.Loan, 0)
	for _, loan := range m.loans {
		if userID != nil && loan.UserID != *userID {
			continue
		}
		if status != nil && loan.Status != *status {
			continue
		}
		loans = append(loans, loan)
	}
	return loans, len(loans), nil
}

func (m *mockLoanRepository) ListWithDetails(ctx context.Context, page, limit int, userID *uuid.UUID, status *string) ([]*repository.LoanWithDetails, int, error) {
	loans := make([]*repository.LoanWithDetails, 0)
	for _, loan := range m.loans {
		if userID != nil && loan.UserID != *userID {
			continue
		}
		if status != nil && loan.Status != *status {
			continue
		}
		loans = append(loans, &repository.LoanWithDetails{
			Loan:      loan,
			UserName:  "Test User",
			BookTitle: "Test Book",
		})
	}
	return loans, len(loans), nil
}

func (m *mockLoanRepository) Update(ctx context.Context, loan *entity.Loan) error {
	m.loans[loan.ID] = loan
	return nil
}

func TestLoanUseCase_BorrowBook(t *testing.T) {
	ctx := context.Background()

	createTestData := func() (*loanUseCase, *entity.User, *entity.Book) {
		userRepo := newMockUserRepository()
		bookRepo := newMockBookRepository()
		loanRepo := newMockLoanRepository()

		userUC := NewUserUseCase(userRepo)
		bookUC := NewBookUseCase(bookRepo)

		user, _ := userUC.Create(ctx, CreateUserInput{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		})

		book, _ := bookUC.Create(ctx, CreateBookInput{
			Title:         "Clean Code",
			Author:        "Robert C. Martin",
			ISBN:          "9780132350884",
			PublishedYear: 2008,
			TotalCopies:   3,
		})

		loanUC := NewLoanUseCase(loanRepo, bookRepo, userRepo).(*loanUseCase)

		return loanUC, user, book
	}

	t.Run("borrow available book", func(t *testing.T) {
		loanUC, user, book := createTestData()

		loan, err := loanUC.BorrowBook(ctx, BorrowBookInput{
			UserID: user.ID,
			BookID: book.ID,
		})
		if err != nil {
			t.Errorf("LoanUseCase.BorrowBook() unexpected error = %v", err)
			return
		}

		if loan.Loan.UserID != user.ID {
			t.Errorf("LoanUseCase.BorrowBook() userID = %v, want %v", loan.Loan.UserID, user.ID)
		}
		if loan.Loan.BookID != book.ID {
			t.Errorf("LoanUseCase.BorrowBook() bookID = %v, want %v", loan.Loan.BookID, book.ID)
		}
		if loan.Loan.Status != entity.LoanStatusActive {
			t.Errorf("LoanUseCase.BorrowBook() status = %v, want %v", loan.Loan.Status, entity.LoanStatusActive)
		}
	})

	t.Run("borrow with custom due date", func(t *testing.T) {
		loanUC, user, book := createTestData()

		dueDate := time.Now().AddDate(0, 0, 30)
		loan, err := loanUC.BorrowBook(ctx, BorrowBookInput{
			UserID:  user.ID,
			BookID:  book.ID,
			DueDate: &dueDate,
		})
		if err != nil {
			t.Errorf("LoanUseCase.BorrowBook() unexpected error = %v", err)
			return
		}

		if loan.Loan.DueDate.Day() != dueDate.Day() {
			t.Errorf("LoanUseCase.BorrowBook() dueDate = %v, want %v", loan.Loan.DueDate, dueDate)
		}
	})

	t.Run("borrow non-existing book", func(t *testing.T) {
		loanUC, user, _ := createTestData()

		_, err := loanUC.BorrowBook(ctx, BorrowBookInput{
			UserID: user.ID,
			BookID: uuid.New(),
		})
		if err != entity.ErrBookNotFound {
			t.Errorf("LoanUseCase.BorrowBook() error = %v, wantErr %v", err, entity.ErrBookNotFound)
		}
	})

	t.Run("borrow by non-existing user", func(t *testing.T) {
		loanUC, _, book := createTestData()

		_, err := loanUC.BorrowBook(ctx, BorrowBookInput{
			UserID: uuid.New(),
			BookID: book.ID,
		})
		if err != entity.ErrUserNotFound {
			t.Errorf("LoanUseCase.BorrowBook() error = %v, wantErr %v", err, entity.ErrUserNotFound)
		}
	})

	t.Run("borrow by disabled user", func(t *testing.T) {
		userRepo := newMockUserRepository()
		bookRepo := newMockBookRepository()
		loanRepo := newMockLoanRepository()

		userUC := NewUserUseCase(userRepo)
		bookUC := NewBookUseCase(bookRepo)

		user, _ := userUC.Create(ctx, CreateUserInput{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		})
		_ = userUC.Disable(ctx, user.ID)

		book, _ := bookUC.Create(ctx, CreateBookInput{
			Title:         "Clean Code",
			Author:        "Robert C. Martin",
			ISBN:          "9780132350884",
			PublishedYear: 2008,
			TotalCopies:   3,
		})

		loanUC := NewLoanUseCase(loanRepo, bookRepo, userRepo)

		_, err := loanUC.BorrowBook(ctx, BorrowBookInput{
			UserID: user.ID,
			BookID: book.ID,
		})
		if err != entity.ErrUserDisabled {
			t.Errorf("LoanUseCase.BorrowBook() error = %v, wantErr %v", err, entity.ErrUserDisabled)
		}
	})
}

func TestLoanUseCase_ReturnBook(t *testing.T) {
	ctx := context.Background()

	t.Run("return borrowed book", func(t *testing.T) {
		userRepo := newMockUserRepository()
		bookRepo := newMockBookRepository()
		loanRepo := newMockLoanRepository()

		userUC := NewUserUseCase(userRepo)
		bookUC := NewBookUseCase(bookRepo)

		user, _ := userUC.Create(ctx, CreateUserInput{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		})

		book, _ := bookUC.Create(ctx, CreateBookInput{
			Title:         "Clean Code",
			Author:        "Robert C. Martin",
			ISBN:          "9780132350884",
			PublishedYear: 2008,
			TotalCopies:   3,
		})

		loanUC := NewLoanUseCase(loanRepo, bookRepo, userRepo)

		borrowed, _ := loanUC.BorrowBook(ctx, BorrowBookInput{
			UserID: user.ID,
			BookID: book.ID,
		})

		returned, err := loanUC.ReturnBook(ctx, borrowed.Loan.ID)
		if err != nil {
			t.Errorf("LoanUseCase.ReturnBook() unexpected error = %v", err)
			return
		}

		if returned.Loan.Status != entity.LoanStatusReturned {
			t.Errorf("LoanUseCase.ReturnBook() status = %v, want %v", returned.Loan.Status, entity.LoanStatusReturned)
		}
		if returned.Loan.ReturnedAt == nil {
			t.Error("LoanUseCase.ReturnBook() returnedAt should not be nil")
		}
	})

	t.Run("return non-existing loan", func(t *testing.T) {
		loanRepo := newMockLoanRepository()
		bookRepo := newMockBookRepository()
		userRepo := newMockUserRepository()

		loanUC := NewLoanUseCase(loanRepo, bookRepo, userRepo)

		_, err := loanUC.ReturnBook(ctx, uuid.New())
		if err != entity.ErrLoanNotFound {
			t.Errorf("LoanUseCase.ReturnBook() error = %v, wantErr %v", err, entity.ErrLoanNotFound)
		}
	})
}

func TestLoanUseCase_List(t *testing.T) {
	ctx := context.Background()

	userRepo := newMockUserRepository()
	bookRepo := newMockBookRepository()
	loanRepo := newMockLoanRepository()

	userUC := NewUserUseCase(userRepo)
	bookUC := NewBookUseCase(bookRepo)

	user, _ := userUC.Create(ctx, CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	})

	book, _ := bookUC.Create(ctx, CreateBookInput{
		Title:         "Clean Code",
		Author:        "Robert C. Martin",
		ISBN:          "9780132350884",
		PublishedYear: 2008,
		TotalCopies:   3,
	})

	loanUC := NewLoanUseCase(loanRepo, bookRepo, userRepo)

	_, _ = loanUC.BorrowBook(ctx, BorrowBookInput{
		UserID: user.ID,
		BookID: book.ID,
	})

	t.Run("list all loans", func(t *testing.T) {
		loans, total, err := loanUC.List(ctx, 1, 10, nil, nil)
		if err != nil {
			t.Errorf("LoanUseCase.List() unexpected error = %v", err)
			return
		}

		if total != 1 {
			t.Errorf("LoanUseCase.List() total = %v, want %v", total, 1)
		}
		if len(loans) != 1 {
			t.Errorf("LoanUseCase.List() len = %v, want %v", len(loans), 1)
		}
	})

	t.Run("list loans by user", func(t *testing.T) {
		loans, total, err := loanUC.List(ctx, 1, 10, &user.ID, nil)
		if err != nil {
			t.Errorf("LoanUseCase.List() unexpected error = %v", err)
			return
		}

		if total != 1 {
			t.Errorf("LoanUseCase.List() total = %v, want %v", total, 1)
		}
		if len(loans) != 1 {
			t.Errorf("LoanUseCase.List() len = %v, want %v", len(loans), 1)
		}
	})

	t.Run("list loans by status", func(t *testing.T) {
		status := entity.LoanStatusActive
		loans, total, err := loanUC.List(ctx, 1, 10, nil, &status)
		if err != nil {
			t.Errorf("LoanUseCase.List() unexpected error = %v", err)
			return
		}

		if total != 1 {
			t.Errorf("LoanUseCase.List() total = %v, want %v", total, 1)
		}
		if len(loans) != 1 {
			t.Errorf("LoanUseCase.List() len = %v, want %v", len(loans), 1)
		}
	})
}
