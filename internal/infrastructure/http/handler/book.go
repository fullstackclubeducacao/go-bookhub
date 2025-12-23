package handler

import (
	"net/http"

	"bookhub/api/generated"
	"bookhub/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Book handlers

func (h *Handler) ListBooks(c *gin.Context, params generated.ListBooksParams) {
	page := 1
	limit := 10

	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	books, total, err := h.bookUseCase.List(c.Request.Context(), page, limit, params.Available)
	if err != nil {
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Error: strPtr("failed to list books"),
			Code:  strPtr("INTERNAL_ERROR"),
		})
		return
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, generated.BookListResponse{
		Data:       booksToResponse(books),
		Pagination: paginationResponse(page, limit, total, totalPages),
	})
}

func (h *Handler) CreateBook(c *gin.Context) {
	var req generated.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("invalid request body"),
			Code:  strPtr("BAD_REQUEST"),
		})
		return
	}

	publishedYear := 0
	if req.PublishedYear != nil {
		publishedYear = *req.PublishedYear
	}

	book, err := h.bookUseCase.Create(c.Request.Context(), usecase.CreateBookInput{
		Title:         req.Title,
		Author:        req.Author,
		ISBN:          req.Isbn,
		PublishedYear: publishedYear,
		TotalCopies:   req.TotalCopies,
	})
	if err != nil {
		handleBookError(c, err)
		return
	}

	c.JSON(http.StatusCreated, generated.BookResponse{
		Data: bookToResponse(book),
	})
}

func (h *Handler) GetBookById(c *gin.Context, id openapi_types.UUID) {
	bookID, err := uuid.Parse(id.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("invalid book ID"),
			Code:  strPtr("BAD_REQUEST"),
		})
		return
	}

	book, err := h.bookUseCase.GetByID(c.Request.Context(), bookID)
	if err != nil {
		handleBookError(c, err)
		return
	}

	c.JSON(http.StatusOK, generated.BookResponse{
		Data: bookToResponse(book),
	})
}
