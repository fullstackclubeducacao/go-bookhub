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
	"bookhub/internal/domain/repository"
	"bookhub/internal/usecase"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListLoans_Success(t *testing.T) {
	handler, _, _, mockLoanUseCase, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	userID := uuid.New()
	bookID := uuid.New()
	loans := []*repository.LoanWithDetails{
		createTestLoanWithDetails(userID, bookID),
		createTestLoanWithDetails(userID, bookID),
	}

	mockLoanUseCase.EXPECT().
		List(gomock.Any(), 1, 10, (*uuid.UUID)(nil), (*string)(nil)).
		Return(loans, 2, nil)

	req := httptest.NewRequest(http.MethodGet, "/loans", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.LoanListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, *response.Data, 2)
}

func TestListLoans_WithStatusFilter(t *testing.T) {
	handler, _, _, mockLoanUseCase, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	userID := uuid.New()
	bookID := uuid.New()
	loans := []*repository.LoanWithDetails{createTestLoanWithDetails(userID, bookID)}
	status := "active"

	mockLoanUseCase.EXPECT().
		List(gomock.Any(), 1, 10, (*uuid.UUID)(nil), &status).
		Return(loans, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/loans?status=active", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListLoans_Error(t *testing.T) {
	handler, _, _, mockLoanUseCase, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	mockLoanUseCase.EXPECT().
		List(gomock.Any(), 1, 10, (*uuid.UUID)(nil), (*string)(nil)).
		Return(nil, 0, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/loans", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBorrowBook_Success(t *testing.T) {
	handler, _, _, mockLoanUseCase, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	userID := uuid.New()
	bookID := uuid.New()
	loanWithDetails := createTestLoanWithDetails(userID, bookID)

	mockLoanUseCase.EXPECT().
		BorrowBook(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx interface{}, input usecase.BorrowBookInput) (*repository.LoanWithDetails, error) {
			if input.UserID == userID && input.BookID == bookID {
				return loanWithDetails, nil
			}
			return nil, errors.New("unexpected input")
		})

	reqBody := generated.BorrowBookRequest{
		UserId: openapi_types.UUID(userID),
		BookId: openapi_types.UUID(bookID),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/loans/borrow", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response generated.LoanResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
}

func TestBorrowBook_BookNotAvailable(t *testing.T) {
	handler, _, _, mockLoanUseCase, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	userID := uuid.New()
	bookID := uuid.New()

	mockLoanUseCase.EXPECT().
		BorrowBook(gomock.Any(), gomock.Any()).
		Return(nil, entity.ErrBookNotAvailable)

	reqBody := generated.BorrowBookRequest{
		UserId: openapi_types.UUID(userID),
		BookId: openapi_types.UUID(bookID),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/loans/borrow", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "book is not available - all copies are borrowed", *response.Error)
}

func TestBorrowBook_UserHasActiveLoan(t *testing.T) {
	handler, _, _, mockLoanUseCase, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	userID := uuid.New()
	bookID := uuid.New()

	mockLoanUseCase.EXPECT().
		BorrowBook(gomock.Any(), gomock.Any()).
		Return(nil, entity.ErrUserHasActiveLoan)

	reqBody := generated.BorrowBookRequest{
		UserId: openapi_types.UUID(userID),
		BookId: openapi_types.UUID(bookID),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/loans/borrow", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user already has an active loan for this book", *response.Error)
}

func TestBorrowBook_UserNotFound(t *testing.T) {
	handler, _, _, mockLoanUseCase, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	userID := uuid.New()
	bookID := uuid.New()

	mockLoanUseCase.EXPECT().
		BorrowBook(gomock.Any(), gomock.Any()).
		Return(nil, entity.ErrUserNotFound)

	reqBody := generated.BorrowBookRequest{
		UserId: openapi_types.UUID(userID),
		BookId: openapi_types.UUID(bookID),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/loans/borrow", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBorrowBook_InvalidRequestBody(t *testing.T) {
	handler, _, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/loans/borrow", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReturnBook_Success(t *testing.T) {
	handler, _, _, mockLoanUseCase, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	userID := uuid.New()
	bookID := uuid.New()
	loanWithDetails := createTestLoanWithDetails(userID, bookID)
	loanID := loanWithDetails.Loan.ID

	mockLoanUseCase.EXPECT().
		ReturnBook(gomock.Any(), loanID).
		Return(loanWithDetails, nil)

	req := httptest.NewRequest(http.MethodPatch, "/loans/"+loanID.String()+"/return", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.LoanResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
}

func TestReturnBook_NotFound(t *testing.T) {
	handler, _, _, mockLoanUseCase, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	loanID := uuid.New()

	mockLoanUseCase.EXPECT().
		ReturnBook(gomock.Any(), loanID).
		Return(nil, entity.ErrLoanNotFound)

	req := httptest.NewRequest(http.MethodPatch, "/loans/"+loanID.String()+"/return", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestReturnBook_AlreadyReturned(t *testing.T) {
	handler, _, _, mockLoanUseCase, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	loanID := uuid.New()

	mockLoanUseCase.EXPECT().
		ReturnBook(gomock.Any(), loanID).
		Return(nil, entity.ErrLoanAlreadyReturned)

	req := httptest.NewRequest(http.MethodPatch, "/loans/"+loanID.String()+"/return", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "loan has already been returned", *response.Error)
}

func TestReturnBook_InvalidID(t *testing.T) {
	handler, _, _, _, _, ctrl := setupTestHandler(t)
	defer ctrl.Finish()
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodPatch, "/loans/invalid-uuid/return", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
