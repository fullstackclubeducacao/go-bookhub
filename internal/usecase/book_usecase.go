package usecase

import (
	"context"

	"bookhub/internal/domain/entity"
	"bookhub/internal/domain/repository"

	"github.com/google/uuid"
)

type BookUseCase interface {
	Create(ctx context.Context, input CreateBookInput) (*entity.Book, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Book, error)
	List(ctx context.Context, page, limit int, availableOnly *bool) ([]*entity.Book, int, error)
}

type CreateBookInput struct {
	Title         string
	Author        string
	ISBN          string
	PublishedYear int
	TotalCopies   int
}

type bookUseCase struct {
	bookRepo repository.BookRepository
}

func NewBookUseCase(bookRepo repository.BookRepository) BookUseCase {
	return &bookUseCase{
		bookRepo: bookRepo,
	}
}

func (uc *bookUseCase) Create(ctx context.Context, input CreateBookInput) (*entity.Book, error) {
	existingBook, err := uc.bookRepo.GetByISBN(ctx, input.ISBN)
	if err == nil && existingBook != nil {
		return nil, entity.ErrInvalidBookISBN
	}

	book, err := entity.NewBook(input.Title, input.Author, input.ISBN, input.PublishedYear, input.TotalCopies)
	if err != nil {
		return nil, err
	}

	if err := uc.bookRepo.Create(ctx, book); err != nil {
		return nil, err
	}

	return book, nil
}

func (uc *bookUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Book, error) {
	book, err := uc.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, entity.ErrBookNotFound
	}
	return book, nil
}

func (uc *bookUseCase) List(ctx context.Context, page, limit int, availableOnly *bool) ([]*entity.Book, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return uc.bookRepo.List(ctx, page, limit, availableOnly)
}
