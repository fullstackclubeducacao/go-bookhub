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

func TestListBooks_Success(t *testing.T) {
	handler, _, mockBookUseCase, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	books := []*entity.Book{createTestBook(), createTestBook()}

	mockBookUseCase.EXPECT().
		List(gomock.Any(), 1, 10, (*bool)(nil)).
		Return(books, 2, nil)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.BookListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, *response.Data, 2)
}

func TestListBooks_WithAvailableFilter(t *testing.T) {
	handler, _, mockBookUseCase, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	books := []*entity.Book{createTestBook()}
	available := true

	mockBookUseCase.EXPECT().
		List(gomock.Any(), 1, 10, &available).
		Return(books, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/books?available=true", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListBooks_Error(t *testing.T) {
	handler, _, mockBookUseCase, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	mockBookUseCase.EXPECT().
		List(gomock.Any(), 1, 10, (*bool)(nil)).
		Return(nil, 0, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateBook_Success(t *testing.T) {
	handler, _, mockBookUseCase, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	book := createTestBook()
	publishedYear := 2024

	mockBookUseCase.EXPECT().
		Create(gomock.Any(), usecase.CreateBookInput{
			Title:         "New Book",
			Author:        "New Author",
			ISBN:          "9876543210",
			PublishedYear: 2024,
			TotalCopies:   3,
		}).
		Return(book, nil)

	reqBody := generated.CreateBookRequest{
		Title:         "New Book",
		Author:        "New Author",
		Isbn:          "9876543210",
		PublishedYear: &publishedYear,
		TotalCopies:   3,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response generated.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
}

func TestCreateBook_InvalidISBN(t *testing.T) {
	handler, _, mockBookUseCase, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	mockBookUseCase.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil, entity.ErrInvalidBookISBN)

	reqBody := generated.CreateBookRequest{
		Title:       "New Book",
		Author:      "New Author",
		Isbn:        "invalid",
		TotalCopies: 3,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetBookById_Success(t *testing.T) {
	handler, _, mockBookUseCase, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	book := createTestBook()

	mockBookUseCase.EXPECT().
		GetByID(gomock.Any(), book.ID).
		Return(book, nil)

	req := httptest.NewRequest(http.MethodGet, "/books/"+book.ID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
}

func TestGetBookById_NotFound(t *testing.T) {
	handler, _, mockBookUseCase, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	bookID := uuid.New()

	mockBookUseCase.EXPECT().
		GetByID(gomock.Any(), bookID).
		Return(nil, entity.ErrBookNotFound)

	req := httptest.NewRequest(http.MethodGet, "/books/"+bookID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
