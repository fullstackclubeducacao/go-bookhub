package entity

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidUserName     = errors.New("invalid user name: must be between 3 and 100 characters")
	ErrInvalidUserEmail    = errors.New("invalid user email format")
	ErrInvalidUserPassword = errors.New("invalid password: must be at least 6 characters")
	ErrUserDisabled        = errors.New("user is disabled")
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyExists  = errors.New("email already exists")
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(name, email, passwordHash string) (*User, error) {
	user := &User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) Validate() error {
	if len(u.Name) < 3 || len(u.Name) > 100 {
		return ErrInvalidUserName
	}

	if !isValidEmail(u.Email) {
		return ErrInvalidUserEmail
	}

	if len(u.PasswordHash) < 6 {
		return ErrInvalidUserPassword
	}

	return nil
}

func (u *User) Update(name, email string) error {
	if name != "" {
		if len(name) < 3 || len(name) > 100 {
			return ErrInvalidUserName
		}
		u.Name = name
	}

	if email != "" {
		if !isValidEmail(email) {
			return ErrInvalidUserEmail
		}
		u.Email = email
	}

	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) Disable() error {
	if !u.Active {
		return ErrUserDisabled
	}
	u.Active = false
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) IsActive() bool {
	return u.Active
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
