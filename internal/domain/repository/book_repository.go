package repository

import (
	"context"

	"bookhub/internal/domain/entity"

	"github.com/google/uuid"
)

type BookRepository interface {
	Create(ctx context.Context, book *entity.Book) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Book, error)
	GetByISBN(ctx context.Context, isbn string) (*entity.Book, error)
	List(ctx context.Context, page, limit int, availableOnly *bool) ([]*entity.Book, int, error)
	Update(ctx context.Context, book *entity.Book) error
	Delete(ctx context.Context, id uuid.UUID) error
}
