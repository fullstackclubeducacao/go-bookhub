package usecase

import (
	"context"

	"bookhub/internal/domain/entity"
	"bookhub/internal/domain/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Create(ctx context.Context, input CreateUserInput) (*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	List(ctx context.Context, page, limit int) ([]*entity.User, int, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*entity.User, error)
	Disable(ctx context.Context, id uuid.UUID) error
	ValidateCredentials(ctx context.Context, email, password string) (*entity.User, error)
}

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

type UpdateUserInput struct {
	Name  *string
	Email *string
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (uc *userUseCase) Create(ctx context.Context, input CreateUserInput) (*entity.User, error) {
	existingUser, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, entity.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := entity.NewUser(input.Name, input.Email, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, entity.ErrUserNotFound
	}
	return user, nil
}

func (uc *userUseCase) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, entity.ErrUserNotFound
	}
	return user, nil
}

func (uc *userUseCase) List(ctx context.Context, page, limit int) ([]*entity.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return uc.userRepo.List(ctx, page, limit)
}

func (uc *userUseCase) Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, entity.ErrUserNotFound
	}

	if input.Email != nil && *input.Email != user.Email {
		existingUser, err := uc.userRepo.GetByEmail(ctx, *input.Email)
		if err == nil && existingUser != nil {
			return nil, entity.ErrEmailAlreadyExists
		}
	}

	name := ""
	if input.Name != nil {
		name = *input.Name
	}

	email := ""
	if input.Email != nil {
		email = *input.Email
	}

	if err := user.Update(name, email); err != nil {
		return nil, err
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) Disable(ctx context.Context, id uuid.UUID) error {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return entity.ErrUserNotFound
	}

	if err := user.Disable(); err != nil {
		return err
	}

	return uc.userRepo.Update(ctx, user)
}

func (uc *userUseCase) ValidateCredentials(ctx context.Context, email, password string) (*entity.User, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, entity.ErrUserNotFound
	}

	if !user.Active {
		return nil, entity.ErrUserDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, entity.ErrUserNotFound
	}

	return user, nil
}
