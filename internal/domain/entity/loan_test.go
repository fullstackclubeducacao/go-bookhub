package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewLoan(t *testing.T) {
	userID := uuid.New()
	bookID := uuid.New()

	t.Run("create loan with default due date", func(t *testing.T) {
		loan, err := NewLoan(userID, bookID, nil)
		if err != nil {
			t.Errorf("NewLoan() unexpected error = %v", err)
			return
		}

		if loan.UserID != userID {
			t.Errorf("NewLoan() userID = %v, want %v", loan.UserID, userID)
		}
		if loan.BookID != bookID {
			t.Errorf("NewLoan() bookID = %v, want %v", loan.BookID, bookID)
		}
		if loan.Status != LoanStatusActive {
			t.Errorf("NewLoan() status = %v, want %v", loan.Status, LoanStatusActive)
		}
		if loan.ReturnedAt != nil {
			t.Error("NewLoan() returnedAt should be nil")
		}

		expectedDue := time.Now().AddDate(0, 0, DefaultLoanDays)
		if loan.DueDate.Day() != expectedDue.Day() {
			t.Errorf("NewLoan() dueDate day = %v, want %v", loan.DueDate.Day(), expectedDue.Day())
		}
	})

	t.Run("create loan with custom due date", func(t *testing.T) {
		customDue := time.Now().AddDate(0, 0, 30)
		loan, err := NewLoan(userID, bookID, &customDue)
		if err != nil {
			t.Errorf("NewLoan() unexpected error = %v", err)
			return
		}

		if loan.DueDate.Day() != customDue.Day() {
			t.Errorf("NewLoan() dueDate = %v, want %v", loan.DueDate, customDue)
		}
	})

	t.Run("create loan with past due date", func(t *testing.T) {
		pastDue := time.Now().AddDate(0, 0, -1)
		_, err := NewLoan(userID, bookID, &pastDue)
		if err != ErrInvalidLoanDueDate {
			t.Errorf("NewLoan() error = %v, wantErr %v", err, ErrInvalidLoanDueDate)
		}
	})
}

func TestLoan_Return(t *testing.T) {
	userID := uuid.New()
	bookID := uuid.New()

	t.Run("return active loan", func(t *testing.T) {
		loan, _ := NewLoan(userID, bookID, nil)

		err := loan.Return()
		if err != nil {
			t.Errorf("Loan.Return() unexpected error = %v", err)
		}

		if loan.Status != LoanStatusReturned {
			t.Errorf("Loan.Return() status = %v, want %v", loan.Status, LoanStatusReturned)
		}
		if loan.ReturnedAt == nil {
			t.Error("Loan.Return() returnedAt should not be nil")
		}
	})

	t.Run("return already returned loan", func(t *testing.T) {
		loan, _ := NewLoan(userID, bookID, nil)
		_ = loan.Return()

		err := loan.Return()
		if err != ErrLoanAlreadyReturned {
			t.Errorf("Loan.Return() error = %v, wantErr %v", err, ErrLoanAlreadyReturned)
		}
	})
}

func TestLoan_IsActive(t *testing.T) {
	userID := uuid.New()
	bookID := uuid.New()

	t.Run("active loan", func(t *testing.T) {
		loan, _ := NewLoan(userID, bookID, nil)

		if !loan.IsActive() {
			t.Error("Loan.IsActive() = false, want true")
		}
	})

	t.Run("returned loan", func(t *testing.T) {
		loan, _ := NewLoan(userID, bookID, nil)

		_ = loan.Return()

		if loan.IsActive() {
			t.Error("Loan.IsActive() = true, want false")
		}
	})
}

func TestLoan_IsOverdue(t *testing.T) {
	userID := uuid.New()
	bookID := uuid.New()

	t.Run("not overdue", func(t *testing.T) {
		loan, _ := NewLoan(userID, bookID, nil)

		if loan.IsOverdue() {
			t.Error("Loan.IsOverdue() = true, want false")
		}
	})

	t.Run("overdue", func(t *testing.T) {
		loan, _ := NewLoan(userID, bookID, nil)
		loan.DueDate = time.Now().AddDate(0, 0, -1)

		if !loan.IsOverdue() {
			t.Error("Loan.IsOverdue() = false, want true")
		}
	})

	t.Run("returned loan is not overdue", func(t *testing.T) {
		loan, _ := NewLoan(userID, bookID, nil)
		loan.DueDate = time.Now().AddDate(0, 0, -1)
		_ = loan.Return()

		if loan.IsOverdue() {
			t.Error("Loan.IsOverdue() returned loan should not be overdue")
		}
	})
}
