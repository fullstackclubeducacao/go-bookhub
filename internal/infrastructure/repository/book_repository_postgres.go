package repository

import (
	"context"
	"database/sql"

	"bookhub/internal/domain/entity"
	"bookhub/internal/domain/repository"
	"bookhub/internal/infrastructure/database/sqlc"

	"github.com/google/uuid"
)

type postgresBookRepository struct {
	queries *sqlc.Queries
}

func NewPostgresBookRepository(db *sql.DB) repository.BookRepository {
	return &postgresBookRepository{
		queries: sqlc.New(db),
	}
}

func (r *postgresBookRepository) Create(ctx context.Context, book *entity.Book) error {
	_, err := r.queries.CreateBook(ctx, sqlc.CreateBookParams{
		ID:              book.ID,
		Title:           book.Title,
		Author:          book.Author,
		Isbn:            book.ISBN,
		PublishedYear:   sql.NullInt32{Int32: int32(book.PublishedYear), Valid: book.PublishedYear > 0},
		TotalCopies:     int32(book.TotalCopies),
		AvailableCopies: int32(book.AvailableCopies),
		CreatedAt:       book.CreatedAt,
		UpdatedAt:       book.UpdatedAt,
	})
	return err
}

func (r *postgresBookRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Book, error) {
	row, err := r.queries.GetBookByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(row), nil
}

func (r *postgresBookRepository) GetByISBN(ctx context.Context, isbn string) (*entity.Book, error) {
	row, err := r.queries.GetBookByISBN(ctx, isbn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(row), nil
}

func (r *postgresBookRepository) List(ctx context.Context, page, limit int, availableOnly *bool) ([]*entity.Book, int, error) {
	offset := (page - 1) * limit

	var rows []sqlc.Book
	var count int64
	var err error

	if availableOnly != nil && *availableOnly {
		rows, err = r.queries.ListAvailableBooks(ctx, sqlc.ListAvailableBooksParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			return nil, 0, err
		}
		count, err = r.queries.CountAvailableBooks(ctx)
	} else {
		rows, err = r.queries.ListBooks(ctx, sqlc.ListBooksParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			return nil, 0, err
		}
		count, err = r.queries.CountBooks(ctx)
	}

	if err != nil {
		return nil, 0, err
	}

	books := make([]*entity.Book, len(rows))
	for i, row := range rows {
		books[i] = r.toEntity(row)
	}

	return books, int(count), nil
}

func (r *postgresBookRepository) Update(ctx context.Context, book *entity.Book) error {
	_, err := r.queries.UpdateBook(ctx, sqlc.UpdateBookParams{
		ID:              book.ID,
		Title:           book.Title,
		Author:          book.Author,
		Isbn:            book.ISBN,
		PublishedYear:   sql.NullInt32{Int32: int32(book.PublishedYear), Valid: book.PublishedYear > 0},
		TotalCopies:     int32(book.TotalCopies),
		AvailableCopies: int32(book.AvailableCopies),
		UpdatedAt:       book.UpdatedAt,
	})
	return err
}

func (r *postgresBookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteBook(ctx, id)
}

func (r *postgresBookRepository) toEntity(row sqlc.Book) *entity.Book {
	publishedYear := 0
	if row.PublishedYear.Valid {
		publishedYear = int(row.PublishedYear.Int32)
	}

	return &entity.Book{
		ID:              row.ID,
		Title:           row.Title,
		Author:          row.Author,
		ISBN:            row.Isbn,
		PublishedYear:   publishedYear,
		TotalCopies:     int(row.TotalCopies),
		AvailableCopies: int(row.AvailableCopies),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}
