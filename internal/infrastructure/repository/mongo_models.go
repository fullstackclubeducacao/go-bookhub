package repository

import (
	"time"

	"bookhub/internal/domain/entity"

	"github.com/google/uuid"
)

// MongoDB document models with bson tags

type userDocument struct {
	ID           uuid.UUID `bson:"id"`
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"passwordhash"`
	Active       bool      `bson:"active"`
	CreatedAt    time.Time `bson:"createdat"`
	UpdatedAt    time.Time `bson:"updatedat"`
}

func toUserDocument(u *entity.User) *userDocument {
	return &userDocument{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Active:       u.Active,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func (d *userDocument) toEntity() *entity.User {
	return &entity.User{
		ID:           d.ID,
		Name:         d.Name,
		Email:        d.Email,
		PasswordHash: d.PasswordHash,
		Active:       d.Active,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

type bookDocument struct {
	ID              uuid.UUID `bson:"id"`
	Title           string    `bson:"title"`
	Author          string    `bson:"author"`
	ISBN            string    `bson:"isbn"`
	PublishedYear   int       `bson:"publishedyear"`
	TotalCopies     int       `bson:"totalcopies"`
	AvailableCopies int       `bson:"availablecopies"`
	CreatedAt       time.Time `bson:"createdat"`
	UpdatedAt       time.Time `bson:"updatedat"`
}

func toBookDocument(b *entity.Book) *bookDocument {
	return &bookDocument{
		ID:              b.ID,
		Title:           b.Title,
		Author:          b.Author,
		ISBN:            b.ISBN,
		PublishedYear:   b.PublishedYear,
		TotalCopies:     b.TotalCopies,
		AvailableCopies: b.AvailableCopies,
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       b.UpdatedAt,
	}
}

func (d *bookDocument) toEntity() *entity.Book {
	return &entity.Book{
		ID:              d.ID,
		Title:           d.Title,
		Author:          d.Author,
		ISBN:            d.ISBN,
		PublishedYear:   d.PublishedYear,
		TotalCopies:     d.TotalCopies,
		AvailableCopies: d.AvailableCopies,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

type loanDocument struct {
	ID         uuid.UUID  `bson:"id"`
	UserID     uuid.UUID  `bson:"userid"`
	BookID     uuid.UUID  `bson:"bookid"`
	BorrowedAt time.Time  `bson:"borrowedat"`
	DueDate    time.Time  `bson:"duedate"`
	ReturnedAt *time.Time `bson:"returnedat"`
	Status     string     `bson:"status"`
}

func toLoanDocument(l *entity.Loan) *loanDocument {
	return &loanDocument{
		ID:         l.ID,
		UserID:     l.UserID,
		BookID:     l.BookID,
		BorrowedAt: l.BorrowedAt,
		DueDate:    l.DueDate,
		ReturnedAt: l.ReturnedAt,
		Status:     l.Status,
	}
}

func (d *loanDocument) toEntity() *entity.Loan {
	return &entity.Loan{
		ID:         d.ID,
		UserID:     d.UserID,
		BookID:     d.BookID,
		BorrowedAt: d.BorrowedAt,
		DueDate:    d.DueDate,
		ReturnedAt: d.ReturnedAt,
		Status:     d.Status,
	}
}
