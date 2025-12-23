package repository

import (
	"context"

	"bookhub/internal/domain/entity"

	"github.com/google/uuid"
)

type LoanRepository interface {
	Create(ctx context.Context, loan *entity.Loan) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Loan, error)
	GetActiveByUserAndBook(ctx context.Context, userID, bookID uuid.UUID) (*entity.Loan, error)
	List(ctx context.Context, page, limit int, userID *uuid.UUID, status *string) ([]*entity.Loan, int, error)
	Update(ctx context.Context, loan *entity.Loan) error
}

type LoanWithDetails struct {
	Loan      *entity.Loan
	UserName  string
	BookTitle string
}

type LoanRepositoryWithDetails interface {
	LoanRepository
	GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*LoanWithDetails, error)
	ListWithDetails(ctx context.Context, page, limit int, userID *uuid.UUID, status *string) ([]*LoanWithDetails, int, error)
}
