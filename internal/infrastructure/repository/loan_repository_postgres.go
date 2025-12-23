package repository

import (
	"context"
	"database/sql"
	"time"

	"bookhub/internal/domain/entity"
	"bookhub/internal/domain/repository"
	"bookhub/internal/infrastructure/database/sqlc"

	"github.com/google/uuid"
)

type postgresLoanRepository struct {
	queries *sqlc.Queries
}

func NewPostgresLoanRepository(db *sql.DB) repository.LoanRepositoryWithDetails {
	return &postgresLoanRepository{
		queries: sqlc.New(db),
	}
}

func (r *postgresLoanRepository) Create(ctx context.Context, loan *entity.Loan) error {
	_, err := r.queries.CreateLoan(ctx, sqlc.CreateLoanParams{
		ID:         loan.ID,
		UserID:     loan.UserID,
		BookID:     loan.BookID,
		BorrowedAt: loan.BorrowedAt,
		DueDate:    loan.DueDate,
		ReturnedAt: r.toNullTime(loan.ReturnedAt),
		Status:     loan.Status,
	})
	return err
}

func (r *postgresLoanRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Loan, error) {
	row, err := r.queries.GetLoanByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(row), nil
}

func (r *postgresLoanRepository) GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*repository.LoanWithDetails, error) {
	row, err := r.queries.GetLoanByIDWithDetails(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntityWithDetails(row), nil
}

func (r *postgresLoanRepository) GetActiveByUserAndBook(ctx context.Context, userID, bookID uuid.UUID) (*entity.Loan, error) {
	row, err := r.queries.GetActiveByUserAndBook(ctx, sqlc.GetActiveByUserAndBookParams{
		UserID: userID,
		BookID: bookID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(row), nil
}

func (r *postgresLoanRepository) List(ctx context.Context, page, limit int, userID *uuid.UUID, status *string) ([]*entity.Loan, int, error) {
	offset := (page - 1) * limit

	var rows []sqlc.Loan
	var count int64
	var err error

	switch {
	case userID != nil && status != nil:
		rows, err = r.queries.ListLoansByUserAndStatus(ctx, sqlc.ListLoansByUserAndStatusParams{
			UserID: *userID,
			Status: *status,
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			return nil, 0, err
		}
		count, err = r.queries.CountLoansByUserAndStatus(ctx, sqlc.CountLoansByUserAndStatusParams{
			UserID: *userID,
			Status: *status,
		})
	case userID != nil:
		rows, err = r.queries.ListLoansByUser(ctx, sqlc.ListLoansByUserParams{
			UserID: *userID,
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			return nil, 0, err
		}
		count, err = r.queries.CountLoansByUser(ctx, *userID)
	case status != nil:
		rows, err = r.queries.ListLoansByStatus(ctx, sqlc.ListLoansByStatusParams{
			Status: *status,
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			return nil, 0, err
		}
		count, err = r.queries.CountLoansByStatus(ctx, *status)
	default:
		rows, err = r.queries.ListLoans(ctx, sqlc.ListLoansParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			return nil, 0, err
		}
		count, err = r.queries.CountLoans(ctx)
	}

	if err != nil {
		return nil, 0, err
	}

	loans := make([]*entity.Loan, len(rows))
	for i, row := range rows {
		loans[i] = r.toEntity(row)
	}

	return loans, int(count), nil
}

func (r *postgresLoanRepository) ListWithDetails(ctx context.Context, page, limit int, userID *uuid.UUID, status *string) ([]*repository.LoanWithDetails, int, error) {
	offset := (page - 1) * limit

	var rows []sqlc.ListLoansWithDetailsRow
	var count int64
	var err error

	switch {
	case userID != nil && status != nil:
		detailRows, e := r.queries.ListLoansByUserAndStatusWithDetails(ctx, sqlc.ListLoansByUserAndStatusWithDetailsParams{
			UserID: *userID,
			Status: *status,
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if e != nil {
			return nil, 0, e
		}
		for _, dr := range detailRows {
			rows = append(rows, sqlc.ListLoansWithDetailsRow(dr))
		}
		count, err = r.queries.CountLoansByUserAndStatus(ctx, sqlc.CountLoansByUserAndStatusParams{
			UserID: *userID,
			Status: *status,
		})
	case userID != nil:
		detailRows, e := r.queries.ListLoansByUserWithDetails(ctx, sqlc.ListLoansByUserWithDetailsParams{
			UserID: *userID,
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if e != nil {
			return nil, 0, e
		}
		for _, dr := range detailRows {
			rows = append(rows, sqlc.ListLoansWithDetailsRow(dr))
		}
		count, err = r.queries.CountLoansByUser(ctx, *userID)
	case status != nil:
		detailRows, e := r.queries.ListLoansByStatusWithDetails(ctx, sqlc.ListLoansByStatusWithDetailsParams{
			Status: *status,
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if e != nil {
			return nil, 0, e
		}
		for _, dr := range detailRows {
			rows = append(rows, sqlc.ListLoansWithDetailsRow(dr))
		}
		count, err = r.queries.CountLoansByStatus(ctx, *status)
	default:
		rows, err = r.queries.ListLoansWithDetails(ctx, sqlc.ListLoansWithDetailsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			return nil, 0, err
		}
		count, err = r.queries.CountLoans(ctx)
	}

	if err != nil {
		return nil, 0, err
	}

	loans := make([]*repository.LoanWithDetails, len(rows))
	for i, row := range rows {
		loans[i] = &repository.LoanWithDetails{
			Loan: &entity.Loan{
				ID:         row.ID,
				UserID:     row.UserID,
				BookID:     row.BookID,
				BorrowedAt: row.BorrowedAt,
				DueDate:    row.DueDate,
				ReturnedAt: r.fromNullTime(row.ReturnedAt),
				Status:     row.Status,
			},
			UserName:  row.UserName,
			BookTitle: row.BookTitle,
		}
	}

	return loans, int(count), nil
}

func (r *postgresLoanRepository) Update(ctx context.Context, loan *entity.Loan) error {
	_, err := r.queries.UpdateLoan(ctx, sqlc.UpdateLoanParams{
		ID:         loan.ID,
		ReturnedAt: r.toNullTime(loan.ReturnedAt),
		Status:     loan.Status,
	})
	return err
}

func (r *postgresLoanRepository) toEntity(row sqlc.Loan) *entity.Loan {
	return &entity.Loan{
		ID:         row.ID,
		UserID:     row.UserID,
		BookID:     row.BookID,
		BorrowedAt: row.BorrowedAt,
		DueDate:    row.DueDate,
		ReturnedAt: r.fromNullTime(row.ReturnedAt),
		Status:     row.Status,
	}
}

func (r *postgresLoanRepository) toEntityWithDetails(row sqlc.GetLoanByIDWithDetailsRow) *repository.LoanWithDetails {
	return &repository.LoanWithDetails{
		Loan: &entity.Loan{
			ID:         row.ID,
			UserID:     row.UserID,
			BookID:     row.BookID,
			BorrowedAt: row.BorrowedAt,
			DueDate:    row.DueDate,
			ReturnedAt: r.fromNullTime(row.ReturnedAt),
			Status:     row.Status,
		},
		UserName:  row.UserName,
		BookTitle: row.BookTitle,
	}
}

func (r *postgresLoanRepository) toNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func (r *postgresLoanRepository) fromNullTime(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}
