package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrLoanNotFound       = errors.New("loan not found")
	ErrLoanAlreadyReturned = errors.New("loan already returned")
	ErrInvalidLoanDueDate = errors.New("invalid due date: must be in the future")
	ErrUserHasActiveLoan  = errors.New("user already has an active loan for this book")
)

const (
	LoanStatusActive   = "active"
	LoanStatusReturned = "returned"
)

const DefaultLoanDays = 14

type Loan struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	BookID     uuid.UUID
	BorrowedAt time.Time
	DueDate    time.Time
	ReturnedAt *time.Time
	Status     string
}

func NewLoan(userID, bookID uuid.UUID, dueDate *time.Time) (*Loan, error) {
	now := time.Now()

	var due time.Time
	if dueDate != nil {
		if dueDate.Before(now) {
			return nil, ErrInvalidLoanDueDate
		}
		due = *dueDate
	} else {
		due = now.AddDate(0, 0, DefaultLoanDays)
	}

	loan := &Loan{
		ID:         uuid.New(),
		UserID:     userID,
		BookID:     bookID,
		BorrowedAt: now,
		DueDate:    due,
		ReturnedAt: nil,
		Status:     LoanStatusActive,
	}

	return loan, nil
}

func (l *Loan) Return() error {
	if l.Status == LoanStatusReturned {
		return ErrLoanAlreadyReturned
	}

	now := time.Now()
	l.ReturnedAt = &now
	l.Status = LoanStatusReturned
	return nil
}

func (l *Loan) IsActive() bool {
	return l.Status == LoanStatusActive
}

func (l *Loan) IsOverdue() bool {
	if !l.IsActive() {
		return false
	}
	return time.Now().After(l.DueDate)
}
