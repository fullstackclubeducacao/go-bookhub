package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidBookTitle       = errors.New("invalid book title: must be between 1 and 200 characters")
	ErrInvalidBookAuthor      = errors.New("invalid book author: must be between 1 and 100 characters")
	ErrInvalidBookISBN        = errors.New("invalid ISBN: must be 10 or 13 digits")
	ErrInvalidTotalCopies     = errors.New("invalid total copies: must be at least 1")
	ErrBookNotFound           = errors.New("book not found")
	ErrBookNotAvailable       = errors.New("book not available: all copies are borrowed")
	ErrInvalidAvailableCopies = errors.New("invalid available copies")
)

const (
	StatusAvailable   = "Disponível"
	StatusUnavailable = "Indisponível - todas as cópias emprestadas"
)

type Book struct {
	ID              uuid.UUID
	Title           string
	Author          string
	ISBN            string
	PublishedYear   int
	TotalCopies     int
	AvailableCopies int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewBook(title, author, isbn string, publishedYear, totalCopies int) (*Book, error) {
	book := &Book{
		ID:              uuid.New(),
		Title:           title,
		Author:          author,
		ISBN:            isbn,
		PublishedYear:   publishedYear,
		TotalCopies:     totalCopies,
		AvailableCopies: totalCopies,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := book.Validate(); err != nil {
		return nil, err
	}

	return book, nil
}

func (b *Book) Validate() error {
	if len(b.Title) < 1 || len(b.Title) > 200 {
		return ErrInvalidBookTitle
	}

	if len(b.Author) < 1 || len(b.Author) > 100 {
		return ErrInvalidBookAuthor
	}

	if !isValidISBN(b.ISBN) {
		return ErrInvalidBookISBN
	}

	if b.TotalCopies < 1 {
		return ErrInvalidTotalCopies
	}

	return nil
}

func (b *Book) IsAvailable() bool {
	return b.AvailableCopies > 0
}

func (b *Book) AvailabilityStatus() string {
	if b.IsAvailable() {
		return StatusAvailable
	}
	return StatusUnavailable
}

func (b *Book) BorrowCopy() error {
	if !b.IsAvailable() {
		return ErrBookNotAvailable
	}
	b.AvailableCopies--
	b.UpdatedAt = time.Now()
	return nil
}

func (b *Book) ReturnCopy() error {
	if b.AvailableCopies >= b.TotalCopies {
		return ErrInvalidAvailableCopies
	}
	b.AvailableCopies++
	b.UpdatedAt = time.Now()
	return nil
}

func isValidISBN(isbn string) bool {
	if len(isbn) != 10 && len(isbn) != 13 {
		return false
	}
	for _, c := range isbn {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
