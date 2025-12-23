package usecase

import (
	"context"
	"testing"

	"bookhub/internal/domain/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	users map[uuid.UUID]*entity.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[uuid.UUID]*entity.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *entity.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) List(ctx context.Context, page, limit int) ([]*entity.User, int, error) {
	users := make([]*entity.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, len(users), nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *entity.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.users, id)
	return nil
}

func TestUserUseCase_Create(t *testing.T) {
	ctx := context.Background()
	repo := newMockUserRepository()
	uc := NewUserUseCase(repo)

	t.Run("create valid user", func(t *testing.T) {
		input := CreateUserInput{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		user, err := uc.Create(ctx, input)
		if err != nil {
			t.Errorf("UserUseCase.Create() unexpected error = %v", err)
			return
		}

		if user.Name != input.Name {
			t.Errorf("UserUseCase.Create() name = %v, want %v", user.Name, input.Name)
		}
		if user.Email != input.Email {
			t.Errorf("UserUseCase.Create() email = %v, want %v", user.Email, input.Email)
		}
		if !user.Active {
			t.Error("UserUseCase.Create() user should be active")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
			t.Error("UserUseCase.Create() password was not hashed correctly")
		}
	})

	t.Run("create user with existing email", func(t *testing.T) {
		input := CreateUserInput{
			Name:     "Jane Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		_, err := uc.Create(ctx, input)
		if err != entity.ErrEmailAlreadyExists {
			t.Errorf("UserUseCase.Create() error = %v, wantErr %v", err, entity.ErrEmailAlreadyExists)
		}
	})

	t.Run("create user with invalid name", func(t *testing.T) {
		input := CreateUserInput{
			Name:     "Jo",
			Email:    "new@example.com",
			Password: "password123",
		}

		_, err := uc.Create(ctx, input)
		if err != entity.ErrInvalidUserName {
			t.Errorf("UserUseCase.Create() error = %v, wantErr %v", err, entity.ErrInvalidUserName)
		}
	})
}

func TestUserUseCase_GetByID(t *testing.T) {
	ctx := context.Background()
	repo := newMockUserRepository()
	uc := NewUserUseCase(repo)

	user, _ := uc.Create(ctx, CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	})

	t.Run("get existing user", func(t *testing.T) {
		found, err := uc.GetByID(ctx, user.ID)
		if err != nil {
			t.Errorf("UserUseCase.GetByID() unexpected error = %v", err)
			return
		}

		if found.ID != user.ID {
			t.Errorf("UserUseCase.GetByID() id = %v, want %v", found.ID, user.ID)
		}
	})

	t.Run("get non-existing user", func(t *testing.T) {
		_, err := uc.GetByID(ctx, uuid.New())
		if err != entity.ErrUserNotFound {
			t.Errorf("UserUseCase.GetByID() error = %v, wantErr %v", err, entity.ErrUserNotFound)
		}
	})
}

func TestUserUseCase_Update(t *testing.T) {
	ctx := context.Background()
	repo := newMockUserRepository()
	uc := NewUserUseCase(repo)

	user, _ := uc.Create(ctx, CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	})

	t.Run("update user name", func(t *testing.T) {
		newName := "John Updated"
		updated, err := uc.Update(ctx, user.ID, UpdateUserInput{
			Name: &newName,
		})
		if err != nil {
			t.Errorf("UserUseCase.Update() unexpected error = %v", err)
			return
		}

		if updated.Name != newName {
			t.Errorf("UserUseCase.Update() name = %v, want %v", updated.Name, newName)
		}
	})

	t.Run("update non-existing user", func(t *testing.T) {
		newName := "Test"
		_, err := uc.Update(ctx, uuid.New(), UpdateUserInput{
			Name: &newName,
		})
		if err != entity.ErrUserNotFound {
			t.Errorf("UserUseCase.Update() error = %v, wantErr %v", err, entity.ErrUserNotFound)
		}
	})
}

func TestUserUseCase_Disable(t *testing.T) {
	ctx := context.Background()
	repo := newMockUserRepository()
	uc := NewUserUseCase(repo)

	user, _ := uc.Create(ctx, CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	})

	t.Run("disable active user", func(t *testing.T) {
		err := uc.Disable(ctx, user.ID)
		if err != nil {
			t.Errorf("UserUseCase.Disable() unexpected error = %v", err)
		}

		found, _ := uc.GetByID(ctx, user.ID)
		if found.Active {
			t.Error("UserUseCase.Disable() user should be disabled")
		}
	})

	t.Run("disable non-existing user", func(t *testing.T) {
		err := uc.Disable(ctx, uuid.New())
		if err != entity.ErrUserNotFound {
			t.Errorf("UserUseCase.Disable() error = %v, wantErr %v", err, entity.ErrUserNotFound)
		}
	})
}

func TestUserUseCase_ValidateCredentials(t *testing.T) {
	ctx := context.Background()
	repo := newMockUserRepository()
	uc := NewUserUseCase(repo)

	_, _ = uc.Create(ctx, CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	})

	t.Run("valid credentials", func(t *testing.T) {
		user, err := uc.ValidateCredentials(ctx, "john@example.com", "password123")
		if err != nil {
			t.Errorf("UserUseCase.ValidateCredentials() unexpected error = %v", err)
			return
		}

		if user.Email != "john@example.com" {
			t.Errorf("UserUseCase.ValidateCredentials() email = %v, want %v", user.Email, "john@example.com")
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		_, err := uc.ValidateCredentials(ctx, "john@example.com", "wrongpassword")
		if err != entity.ErrUserNotFound {
			t.Errorf("UserUseCase.ValidateCredentials() error = %v, wantErr %v", err, entity.ErrUserNotFound)
		}
	})

	t.Run("non-existing email", func(t *testing.T) {
		_, err := uc.ValidateCredentials(ctx, "nonexistent@example.com", "password123")
		if err != entity.ErrUserNotFound {
			t.Errorf("UserUseCase.ValidateCredentials() error = %v, wantErr %v", err, entity.ErrUserNotFound)
		}
	})
}
