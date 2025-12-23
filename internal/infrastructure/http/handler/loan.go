package handler

import (
	"net/http"
	"time"

	"bookhub/api/generated"
	"bookhub/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Loan handlers

func (h *Handler) ListLoans(c *gin.Context, params generated.ListLoansParams) {
	page := 1
	limit := 10

	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	var userID *uuid.UUID
	if params.UserId != nil {
		id, err := uuid.Parse(params.UserId.String())
		if err == nil {
			userID = &id
		}
	}

	var status *string
	if params.Status != nil {
		s := string(*params.Status)
		status = &s
	}

	loans, total, err := h.loanUseCase.List(c.Request.Context(), page, limit, userID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Error: strPtr("failed to list loans"),
			Code:  strPtr("INTERNAL_ERROR"),
		})
		return
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, generated.LoanListResponse{
		Data:       loansToResponse(loans),
		Pagination: paginationResponse(page, limit, total, totalPages),
	})
}

func (h *Handler) BorrowBook(c *gin.Context) {
	var req generated.BorrowBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("invalid request body"),
			Code:  strPtr("BAD_REQUEST"),
		})
		return
	}

	userID, err := uuid.Parse(req.UserId.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("invalid user ID"),
			Code:  strPtr("BAD_REQUEST"),
		})
		return
	}

	bookID, err := uuid.Parse(req.BookId.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("invalid book ID"),
			Code:  strPtr("BAD_REQUEST"),
		})
		return
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		t := req.DueDate.Time
		dueDate = &t
	}

	loan, err := h.loanUseCase.BorrowBook(c.Request.Context(), usecase.BorrowBookInput{
		UserID:  userID,
		BookID:  bookID,
		DueDate: dueDate,
	})
	if err != nil {
		handleLoanError(c, err)
		return
	}

	c.JSON(http.StatusCreated, generated.LoanResponse{
		Data: loanToResponse(loan),
	})
}

func (h *Handler) ReturnBook(c *gin.Context, id openapi_types.UUID) {
	loanID, err := uuid.Parse(id.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("invalid loan ID"),
			Code:  strPtr("BAD_REQUEST"),
		})
		return
	}

	loan, err := h.loanUseCase.ReturnBook(c.Request.Context(), loanID)
	if err != nil {
		handleLoanError(c, err)
		return
	}

	c.JSON(http.StatusOK, generated.LoanResponse{
		Data: loanToResponse(loan),
	})
}
