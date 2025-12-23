package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bookhub/api/generated"
	"bookhub/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLogin_Success(t *testing.T) {
	handler, mockUserUseCase, _, _, mockJWTService, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	user := createTestUser()
	expiresAt := time.Now().Add(24 * time.Hour)

	mockUserUseCase.EXPECT().
		ValidateCredentials(gomock.Any(), "test@example.com", "password123").
		Return(user, nil)
	mockJWTService.EXPECT().
		GenerateToken(user.ID, user.Email).
		Return("test-token", expiresAt, nil)

	reqBody := generated.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test-token", *response.Token)
	assert.NotNil(t, response.User)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	handler, mockUserUseCase, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	mockUserUseCase.EXPECT().
		ValidateCredentials(gomock.Any(), "test@example.com", "wrongpassword").
		Return(nil, entity.ErrUserNotFound)

	reqBody := generated.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid credentials", *response.Error)
}

func TestLogin_InvalidRequestBody(t *testing.T) {
	handler, _, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_TokenGenerationError(t *testing.T) {
	handler, mockUserUseCase, _, _, mockJWTService, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	user := createTestUser()

	mockUserUseCase.EXPECT().
		ValidateCredentials(gomock.Any(), "test@example.com", "password123").
		Return(user, nil)
	mockJWTService.EXPECT().
		GenerateToken(user.ID, user.Email).
		Return("", time.Time{}, errors.New("token generation failed"))

	reqBody := generated.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
