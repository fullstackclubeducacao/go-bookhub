package handler

import (
	"testing"
	"time"

	"bookhub/api/generated"
	"bookhub/internal/domain/entity"
	"bookhub/internal/domain/repository"
	"bookhub/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupTestHandler(t *testing.T) (*Handler, *mocks.MockUserUseCase, *mocks.MockBookUseCase, *mocks.MockLoanUseCase, *mocks.MockJWTService, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	mockBookUseCase := mocks.NewMockBookUseCase(ctrl)
	mockLoanUseCase := mocks.NewMockLoanUseCase(ctrl)
	mockJWTService := mocks.NewMockJWTService(ctrl)

	handler := NewHandler(mockUserUseCase, mockBookUseCase, mockLoanUseCase, mockJWTService)
	return handler, mockUserUseCase, mockBookUseCase, mockLoanUseCase, mockJWTService, ctrl
}

func setupTestRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, handler)
	return router
}

func createTestUser() *entity.User {
	return &entity.User{
		ID:           uuid.New(),
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword123",
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func createTestBook() *entity.Book {
	return &entity.Book{
		ID:              uuid.New(),
		Title:           "Test Book",
		Author:          "Test Author",
		ISBN:            "1234567890",
		PublishedYear:   2024,
		TotalCopies:     5,
		AvailableCopies: 3,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func createTestLoan(userID, bookID uuid.UUID) *entity.Loan {
	now := time.Now()
	return &entity.Loan{
		ID:         uuid.New(),
		UserID:     userID,
		BookID:     bookID,
		BorrowedAt: now,
		DueDate:    now.AddDate(0, 0, 14),
		ReturnedAt: nil,
		Status:     entity.LoanStatusActive,
	}
}

func createTestLoanWithDetails(userID, bookID uuid.UUID) *repository.LoanWithDetails {
	return &repository.LoanWithDetails{
		Loan:      createTestLoan(userID, bookID),
		UserName:  "Test User",
		BookTitle: "Test Book",
	}
}

func TestNewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	mockBookUseCase := mocks.NewMockBookUseCase(ctrl)
	mockLoanUseCase := mocks.NewMockLoanUseCase(ctrl)
	mockJWTService := mocks.NewMockJWTService(ctrl)

	handler := NewHandler(mockUserUseCase, mockBookUseCase, mockLoanUseCase, mockJWTService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockJWTService, handler.JWTService())
}
