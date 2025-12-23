package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"bookhub/api/generated"
	"bookhub/internal/domain/entity"
	"bookhub/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListUsers_Success(t *testing.T) {
	handler, mockUserUseCase, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	users := []*entity.User{createTestUser(), createTestUser()}

	mockUserUseCase.EXPECT().
		List(gomock.Any(), 1, 10).
		Return(users, 2, nil)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.UserListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, *response.Data, 2)
	assert.Equal(t, 2, *response.Pagination.Total)
}

func TestListUsers_WithPagination(t *testing.T) {
	handler, mockUserUseCase, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	users := []*entity.User{createTestUser()}

	mockUserUseCase.EXPECT().
		List(gomock.Any(), 2, 5).
		Return(users, 10, nil)

	req := httptest.NewRequest(http.MethodGet, "/users?page=2&limit=5", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.UserListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, *response.Pagination.Page)
	assert.Equal(t, 5, *response.Pagination.Limit)
	assert.Equal(t, 10, *response.Pagination.Total)
}

func TestListUsers_Error(t *testing.T) {
	handler, mockUserUseCase, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	mockUserUseCase.EXPECT().
		List(gomock.Any(), 1, 10).
		Return(nil, 0, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateUser_Success(t *testing.T) {
	handler, mockUserUseCase, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	user := createTestUser()

	mockUserUseCase.EXPECT().
		Create(gomock.Any(), usecase.CreateUserInput{
			Name:     "New User",
			Email:    "newuser@example.com",
			Password: "password123",
		}).
		Return(user, nil)

	reqBody := generated.CreateUserRequest{
		Name:     "New User",
		Email:    "newuser@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response generated.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	handler, mockUserUseCase, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	mockUserUseCase.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil, entity.ErrEmailAlreadyExists)

	reqBody := generated.CreateUserRequest{
		Name:     "New User",
		Email:    "existing@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "email already exists", *response.Error)
}

func TestCreateUser_InvalidRequestBody(t *testing.T) {
	handler, _, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserById_Success(t *testing.T) {
	handler, mockUserUseCase, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	user := createTestUser()

	mockUserUseCase.EXPECT().
		GetByID(gomock.Any(), user.ID).
		Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/"+user.ID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
}

func TestGetUserById_NotFound(t *testing.T) {
	handler, mockUserUseCase, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	userID := uuid.New()

	mockUserUseCase.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(nil, entity.ErrUserNotFound)

	req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUserById_InvalidID(t *testing.T) {
	handler, _, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/users/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUser_Success(t *testing.T) {
	handler, mockUserUseCase, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	user := createTestUser()
	updatedName := "Updated Name"

	mockUserUseCase.EXPECT().
		Update(gomock.Any(), user.ID, gomock.Any()).
		Return(user, nil)

	reqBody := generated.UpdateUserRequest{
		Name: &updatedName,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/users/"+user.ID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateUser_NotFound(t *testing.T) {
	handler, mockUserUseCase, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	userID := uuid.New()
	updatedName := "Updated Name"

	mockUserUseCase.EXPECT().
		Update(gomock.Any(), userID, gomock.Any()).
		Return(nil, entity.ErrUserNotFound)

	reqBody := generated.UpdateUserRequest{
		Name: &updatedName,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/users/"+userID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDisableUser_Success(t *testing.T) {
	handler, mockUserUseCase, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	userID := uuid.New()

	mockUserUseCase.EXPECT().
		Disable(gomock.Any(), userID).
		Return(nil)

	req := httptest.NewRequest(http.MethodPatch, "/users/"+userID.String()+"/disable", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user disabled successfully", *response.Message)
}

func TestDisableUser_NotFound(t *testing.T) {
	handler, mockUserUseCase, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	userID := uuid.New()

	mockUserUseCase.EXPECT().
		Disable(gomock.Any(), userID).
		Return(entity.ErrUserNotFound)

	req := httptest.NewRequest(http.MethodPatch, "/users/"+userID.String()+"/disable", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDisableUser_AlreadyDisabled(t *testing.T) {
	handler, mockUserUseCase, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	userID := uuid.New()

	mockUserUseCase.EXPECT().
		Disable(gomock.Any(), userID).
		Return(entity.ErrUserDisabled)

	req := httptest.NewRequest(http.MethodPatch, "/users/"+userID.String()+"/disable", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
