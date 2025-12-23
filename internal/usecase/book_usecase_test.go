package usecase

import (
	"context"
	"testing"

	"bookhub/internal/domain/entity"

	"github.com/google/uuid"
)

type mockBookRepository struct {
	books map[uuid.UUID]*entity.Book
}

func newMockBookRepository() *mockBookRepository {
	return &mockBookRepository{
		books: make(map[uuid.UUID]*entity.Book),
	}
}

func (m *mockBookRepository) Create(ctx context.Context, book *entity.Book) error {
	m.books[book.ID] = book
	return nil
}

func (m *mockBookRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Book, error) {
	if book, exists := m.books[id]; exists {
		return book, nil
	}
	return nil, nil
}

func (m *mockBookRepository) GetByISBN(ctx context.Context, isbn string) (*entity.Book, error) {
	for _, book := range m.books {
		if book.ISBN == isbn {
			return book, nil
		}
	}
	return nil, nil
}

func (m *mockBookRepository) List(ctx context.Context, page, limit int, availableOnly *bool) ([]*entity.Book, int, error) {
	books := make([]*entity.Book, 0)
	for _, book := range m.books {
		if availableOnly != nil && *availableOnly && !book.IsAvailable() {
			continue
		}
		books = append(books, book)
	}
	return books, len(books), nil
}

func (m *mockBookRepository) Update(ctx context.Context, book *entity.Book) error {
	m.books[book.ID] = book
	return nil
}

func (m *mockBookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.books, id)
	return nil
}

func TestBookUseCase_Create(t *testing.T) {
	ctx := context.Background()
	repo := newMockBookRepository()
	uc := NewBookUseCase(repo)

	t.Run("create valid book", func(t *testing.T) {
		input := CreateBookInput{
			Title:         "Clean Code",
			Author:        "Robert C. Martin",
			ISBN:          "9780132350884",
			PublishedYear: 2008,
			TotalCopies:   5,
		}

		book, err := uc.Create(ctx, input)
		if err != nil {
			t.Errorf("BookUseCase.Create() unexpected error = %v", err)
			return
		}

		if book.Title != input.Title {
			t.Errorf("BookUseCase.Create() title = %v, want %v", book.Title, input.Title)
		}
		if book.AvailableCopies != input.TotalCopies {
			t.Errorf("BookUseCase.Create() available copies = %v, want %v", book.AvailableCopies, input.TotalCopies)
		}
	})

	t.Run("create book with duplicate ISBN", func(t *testing.T) {
		input := CreateBookInput{
			Title:         "Another Book",
			Author:        "Another Author",
			ISBN:          "9780132350884",
			PublishedYear: 2020,
			TotalCopies:   3,
		}

		_, err := uc.Create(ctx, input)
		if err != entity.ErrInvalidBookISBN {
			t.Errorf("BookUseCase.Create() error = %v, wantErr %v", err, entity.ErrInvalidBookISBN)
		}
	})

	t.Run("create book with invalid title", func(t *testing.T) {
		input := CreateBookInput{
			Title:         "",
			Author:        "Author",
			ISBN:          "9780132350885",
			PublishedYear: 2020,
			TotalCopies:   3,
		}

		_, err := uc.Create(ctx, input)
		if err != entity.ErrInvalidBookTitle {
			t.Errorf("BookUseCase.Create() error = %v, wantErr %v", err, entity.ErrInvalidBookTitle)
		}
	})
}

func TestBookUseCase_GetByID(t *testing.T) {
	ctx := context.Background()
	repo := newMockBookRepository()
	uc := NewBookUseCase(repo)

	book, _ := uc.Create(ctx, CreateBookInput{
		Title:         "Clean Code",
		Author:        "Robert C. Martin",
		ISBN:          "9780132350884",
		PublishedYear: 2008,
		TotalCopies:   5,
	})

	t.Run("get existing book", func(t *testing.T) {
		found, err := uc.GetByID(ctx, book.ID)
		if err != nil {
			t.Errorf("BookUseCase.GetByID() unexpected error = %v", err)
			return
		}

		if found.ID != book.ID {
			t.Errorf("BookUseCase.GetByID() id = %v, want %v", found.ID, book.ID)
		}
	})

	t.Run("get non-existing book", func(t *testing.T) {
		_, err := uc.GetByID(ctx, uuid.New())
		if err != entity.ErrBookNotFound {
			t.Errorf("BookUseCase.GetByID() error = %v, wantErr %v", err, entity.ErrBookNotFound)
		}
	})
}

func TestBookUseCase_List(t *testing.T) {
	ctx := context.Background()
	repo := newMockBookRepository()
	uc := NewBookUseCase(repo)

	_, _ = uc.Create(ctx, CreateBookInput{
		Title:         "Book 1",
		Author:        "Author 1",
		ISBN:          "9780132350881",
		PublishedYear: 2020,
		TotalCopies:   5,
	})

	book2, _ := uc.Create(ctx, CreateBookInput{
		Title:         "Book 2",
		Author:        "Author 2",
		ISBN:          "9780132350882",
		PublishedYear: 2021,
		TotalCopies:   1,
	})
	book2.AvailableCopies = 0
	//nolint
	repo.Update(ctx, book2)

	t.Run("list all books", func(t *testing.T) {
		books, total, err := uc.List(ctx, 1, 10, nil)
		if err != nil {
			t.Errorf("BookUseCase.List() unexpected error = %v", err)
			return
		}

		if total != 2 {
			t.Errorf("BookUseCase.List() total = %v, want %v", total, 2)
		}
		if len(books) != 2 {
			t.Errorf("BookUseCase.List() len = %v, want %v", len(books), 2)
		}
	})

	t.Run("list only available books", func(t *testing.T) {
		availableOnly := true
		books, total, err := uc.List(ctx, 1, 10, &availableOnly)
		if err != nil {
			t.Errorf("BookUseCase.List() unexpected error = %v", err)
			return
		}

		if total != 1 {
			t.Errorf("BookUseCase.List() total = %v, want %v", total, 1)
		}
		if len(books) != 1 {
			t.Errorf("BookUseCase.List() len = %v, want %v", len(books), 1)
		}
	})
}
