package usecase

import (
	"context"
	"time"

	"bookhub/internal/domain/entity"
	"bookhub/internal/domain/repository"

	"github.com/google/uuid"
)

type LoanUseCase interface {
	BorrowBook(ctx context.Context, input BorrowBookInput) (*repository.LoanWithDetails, error)
	ReturnBook(ctx context.Context, loanID uuid.UUID) (*repository.LoanWithDetails, error)
	GetByID(ctx context.Context, id uuid.UUID) (*repository.LoanWithDetails, error)
	List(ctx context.Context, page, limit int, userID *uuid.UUID, status *string) ([]*repository.LoanWithDetails, int, error)
}

type BorrowBookInput struct {
	UserID  uuid.UUID
	BookID  uuid.UUID
	DueDate *time.Time
}

type loanUseCase struct {
	loanRepo repository.LoanRepositoryWithDetails
	bookRepo repository.BookRepository
	userRepo repository.UserRepository
}

func NewLoanUseCase(
	loanRepo repository.LoanRepositoryWithDetails,
	bookRepo repository.BookRepository,
	userRepo repository.UserRepository,
) LoanUseCase {
	return &loanUseCase{
		loanRepo: loanRepo,
		bookRepo: bookRepo,
		userRepo: userRepo,
	}
}

func (uc *loanUseCase) BorrowBook(ctx context.Context, input BorrowBookInput) (*repository.LoanWithDetails, error) {
	user, err := uc.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, entity.ErrUserNotFound
	}
	if !user.Active {
		return nil, entity.ErrUserDisabled
	}

	book, err := uc.bookRepo.GetByID(ctx, input.BookID)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, entity.ErrBookNotFound
	}

	if !book.IsAvailable() {
		return nil, entity.ErrBookNotAvailable
	}

	existingLoan, _ := uc.loanRepo.GetActiveByUserAndBook(ctx, input.UserID, input.BookID)
	if existingLoan != nil {
		return nil, entity.ErrUserHasActiveLoan
	}

	loan, err := entity.NewLoan(input.UserID, input.BookID, input.DueDate)
	if err != nil {
		return nil, err
	}

	if err := book.BorrowCopy(); err != nil {
		return nil, err
	}

	if err := uc.bookRepo.Update(ctx, book); err != nil {
		return nil, err
	}

	if err := uc.loanRepo.Create(ctx, loan); err != nil {
		return nil, err
	}

	return &repository.LoanWithDetails{
		Loan:      loan,
		UserName:  user.Name,
		BookTitle: book.Title,
	}, nil
}

func (uc *loanUseCase) ReturnBook(ctx context.Context, loanID uuid.UUID) (*repository.LoanWithDetails, error) {
	loanDetails, err := uc.loanRepo.GetByIDWithDetails(ctx, loanID)
	if err != nil {
		return nil, err
	}
	if loanDetails == nil || loanDetails.Loan == nil {
		return nil, entity.ErrLoanNotFound
	}

	if err := loanDetails.Loan.Return(); err != nil {
		return nil, err
	}

	book, err := uc.bookRepo.GetByID(ctx, loanDetails.Loan.BookID)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, entity.ErrBookNotFound
	}

	if err := book.ReturnCopy(); err != nil {
		return nil, err
	}

	if err := uc.bookRepo.Update(ctx, book); err != nil {
		return nil, err
	}

	if err := uc.loanRepo.Update(ctx, loanDetails.Loan); err != nil {
		return nil, err
	}

	return loanDetails, nil
}

func (uc *loanUseCase) GetByID(ctx context.Context, id uuid.UUID) (*repository.LoanWithDetails, error) {
	loanDetails, err := uc.loanRepo.GetByIDWithDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if loanDetails == nil {
		return nil, entity.ErrLoanNotFound
	}
	return loanDetails, nil
}

func (uc *loanUseCase) List(ctx context.Context, page, limit int, userID *uuid.UUID, status *string) ([]*repository.LoanWithDetails, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return uc.loanRepo.ListWithDetails(ctx, page, limit, userID, status)
}
