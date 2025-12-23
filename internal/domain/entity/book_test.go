package entity

import (
	"testing"
)

func TestNewBook(t *testing.T) {
	tests := []struct {
		name          string
		title         string
		author        string
		isbn          string
		publishedYear int
		totalCopies   int
		wantErr       error
	}{
		{
			name:          "valid book",
			title:         "Clean Code",
			author:        "Robert C. Martin",
			isbn:          "9780132350884",
			publishedYear: 2008,
			totalCopies:   5,
			wantErr:       nil,
		},
		{
			name:          "valid book with 10 digit ISBN",
			title:         "Clean Code",
			author:        "Robert C. Martin",
			isbn:          "0132350882",
			publishedYear: 2008,
			totalCopies:   3,
			wantErr:       nil,
		},
		{
			name:          "empty title",
			title:         "",
			author:        "Robert C. Martin",
			isbn:          "9780132350884",
			publishedYear: 2008,
			totalCopies:   5,
			wantErr:       ErrInvalidBookTitle,
		},
		{
			name:          "title too long",
			title:         string(make([]byte, 201)),
			author:        "Robert C. Martin",
			isbn:          "9780132350884",
			publishedYear: 2008,
			totalCopies:   5,
			wantErr:       ErrInvalidBookTitle,
		},
		{
			name:          "empty author",
			title:         "Clean Code",
			author:        "",
			isbn:          "9780132350884",
			publishedYear: 2008,
			totalCopies:   5,
			wantErr:       ErrInvalidBookAuthor,
		},
		{
			name:          "invalid ISBN - too short",
			title:         "Clean Code",
			author:        "Robert C. Martin",
			isbn:          "123456789",
			publishedYear: 2008,
			totalCopies:   5,
			wantErr:       ErrInvalidBookISBN,
		},
		{
			name:          "invalid ISBN - contains letters",
			title:         "Clean Code",
			author:        "Robert C. Martin",
			isbn:          "978013235088X",
			publishedYear: 2008,
			totalCopies:   5,
			wantErr:       ErrInvalidBookISBN,
		},
		{
			name:          "zero copies",
			title:         "Clean Code",
			author:        "Robert C. Martin",
			isbn:          "9780132350884",
			publishedYear: 2008,
			totalCopies:   0,
			wantErr:       ErrInvalidTotalCopies,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, err := NewBook(tt.title, tt.author, tt.isbn, tt.publishedYear, tt.totalCopies)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewBook() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewBook() unexpected error = %v", err)
				return
			}

			if book.Title != tt.title {
				t.Errorf("NewBook() title = %v, want %v", book.Title, tt.title)
			}
			if book.AvailableCopies != tt.totalCopies {
				t.Errorf("NewBook() available copies = %v, want %v", book.AvailableCopies, tt.totalCopies)
			}
		})
	}
}

func TestBook_IsAvailable(t *testing.T) {
	tests := []struct {
		name            string
		availableCopies int
		want            bool
	}{
		{
			name:            "has available copies",
			availableCopies: 3,
			want:            true,
		},
		{
			name:            "no available copies",
			availableCopies: 0,
			want:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, _ := NewBook("Test Book", "Author", "9780132350884", 2020, 5)
			book.AvailableCopies = tt.availableCopies

			if got := book.IsAvailable(); got != tt.want {
				t.Errorf("Book.IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBook_AvailabilityStatus(t *testing.T) {
	tests := []struct {
		name            string
		availableCopies int
		wantStatus      string
	}{
		{
			name:            "available",
			availableCopies: 3,
			wantStatus:      StatusAvailable,
		},
		{
			name:            "unavailable",
			availableCopies: 0,
			wantStatus:      StatusUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, _ := NewBook("Test Book", "Author", "9780132350884", 2020, 5)
			book.AvailableCopies = tt.availableCopies

			if got := book.AvailabilityStatus(); got != tt.wantStatus {
				t.Errorf("Book.AvailabilityStatus() = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestBook_BorrowCopy(t *testing.T) {
	t.Run("borrow from available book", func(t *testing.T) {
		book, _ := NewBook("Test Book", "Author", "9780132350884", 2020, 3)
		initialCopies := book.AvailableCopies

		err := book.BorrowCopy()
		if err != nil {
			t.Errorf("Book.BorrowCopy() unexpected error = %v", err)
		}

		if book.AvailableCopies != initialCopies-1 {
			t.Errorf("Book.BorrowCopy() available copies = %v, want %v", book.AvailableCopies, initialCopies-1)
		}
	})

	t.Run("borrow from unavailable book", func(t *testing.T) {
		book, _ := NewBook("Test Book", "Author", "9780132350884", 2020, 1)
		book.AvailableCopies = 0

		err := book.BorrowCopy()
		if err != ErrBookNotAvailable {
			t.Errorf("Book.BorrowCopy() error = %v, wantErr %v", err, ErrBookNotAvailable)
		}
	})
}

func TestBook_ReturnCopy(t *testing.T) {
	t.Run("return borrowed copy", func(t *testing.T) {
		book, _ := NewBook("Test Book", "Author", "9780132350884", 2020, 3)
		book.AvailableCopies = 2

		err := book.ReturnCopy()
		if err != nil {
			t.Errorf("Book.ReturnCopy() unexpected error = %v", err)
		}

		if book.AvailableCopies != 3 {
			t.Errorf("Book.ReturnCopy() available copies = %v, want %v", book.AvailableCopies, 3)
		}
	})

	t.Run("return when all copies available", func(t *testing.T) {
		book, _ := NewBook("Test Book", "Author", "9780132350884", 2020, 3)

		err := book.ReturnCopy()
		if err != ErrInvalidAvailableCopies {
			t.Errorf("Book.ReturnCopy() error = %v, wantErr %v", err, ErrInvalidAvailableCopies)
		}
	})
}

func TestIsValidISBN(t *testing.T) {
	tests := []struct {
		isbn  string
		valid bool
	}{
		{"9780132350884", true},
		{"0132350882", true},
		{"123456789", false},
		{"12345678901234", false},
		{"978013235088X", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.isbn, func(t *testing.T) {
			if got := isValidISBN(tt.isbn); got != tt.valid {
				t.Errorf("isValidISBN(%q) = %v, want %v", tt.isbn, got, tt.valid)
			}
		})
	}
}
