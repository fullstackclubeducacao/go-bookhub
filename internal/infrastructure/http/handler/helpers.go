package handler

import (
	"net/http"
	"time"

	"bookhub/api/generated"
	"bookhub/internal/domain/entity"
	"bookhub/internal/domain/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func strPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func uuidToOpenAPI(id uuid.UUID) *openapi_types.UUID {
	oaUUID := openapi_types.UUID(id)
	return &oaUUID
}

func emailToOpenAPI(email string) *openapi_types.Email {
	oaEmail := openapi_types.Email(email)
	return &oaEmail
}

func userToResponse(user *entity.User) *generated.User {
	if user == nil {
		return nil
	}
	return &generated.User{
		Id:        uuidToOpenAPI(user.ID),
		Name:      &user.Name,
		Email:     emailToOpenAPI(user.Email),
		Active:    &user.Active,
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
	}
}

func usersToResponse(users []*entity.User) *[]generated.User {
	result := make([]generated.User, len(users))
	for i, user := range users {
		u := userToResponse(user)
		if u != nil {
			result[i] = *u
		}
	}
	return &result
}

func bookToResponse(book *entity.Book) *generated.Book {
	if book == nil {
		return nil
	}
	status := book.AvailabilityStatus()
	return &generated.Book{
		Id:                 uuidToOpenAPI(book.ID),
		Title:              &book.Title,
		Author:             &book.Author,
		Isbn:               &book.ISBN,
		PublishedYear:      &book.PublishedYear,
		TotalCopies:        &book.TotalCopies,
		AvailableCopies:    &book.AvailableCopies,
		AvailabilityStatus: &status,
		CreatedAt:          &book.CreatedAt,
		UpdatedAt:          &book.UpdatedAt,
	}
}

func booksToResponse(books []*entity.Book) *[]generated.Book {
	result := make([]generated.Book, len(books))
	for i, book := range books {
		b := bookToResponse(book)
		if b != nil {
			result[i] = *b
		}
	}
	return &result
}

func loanToResponse(loan *repository.LoanWithDetails) *generated.Loan {
	if loan == nil || loan.Loan == nil {
		return nil
	}
	status := generated.LoanStatus(loan.Loan.Status)

	result := &generated.Loan{
		Id:         uuidToOpenAPI(loan.Loan.ID),
		UserId:     uuidToOpenAPI(loan.Loan.UserID),
		UserName:   &loan.UserName,
		BookId:     uuidToOpenAPI(loan.Loan.BookID),
		BookTitle:  &loan.BookTitle,
		BorrowedAt: &loan.Loan.BorrowedAt,
		DueDate:    &loan.Loan.DueDate,
		Status:     &status,
	}

	if loan.Loan.ReturnedAt != nil {
		result.ReturnedAt = loan.Loan.ReturnedAt
	}

	return result
}

func loansToResponse(loans []*repository.LoanWithDetails) *[]generated.Loan {
	result := make([]generated.Loan, len(loans))
	for i, loan := range loans {
		l := loanToResponse(loan)
		if l != nil {
			result[i] = *l
		}
	}
	return &result
}

func paginationResponse(page, limit, total, totalPages int) *generated.Pagination {
	return &generated.Pagination{
		Page:       &page,
		Limit:      &limit,
		Total:      &total,
		TotalPages: &totalPages,
	}
}

func handleUserError(c *gin.Context, err error) {
	switch err {
	case entity.ErrUserNotFound:
		c.JSON(http.StatusNotFound, generated.ErrorResponse{
			Error: strPtr("user not found"),
			Code:  strPtr("NOT_FOUND"),
		})
	case entity.ErrEmailAlreadyExists:
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("email already exists"),
			Code:  strPtr("EMAIL_EXISTS"),
		})
	case entity.ErrUserDisabled:
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("user is already disabled"),
			Code:  strPtr("USER_DISABLED"),
		})
	case entity.ErrInvalidUserName, entity.ErrInvalidUserEmail, entity.ErrInvalidUserPassword:
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr(err.Error()),
			Code:  strPtr("VALIDATION_ERROR"),
		})
	default:
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Error: strPtr("internal server error"),
			Code:  strPtr("INTERNAL_ERROR"),
		})
	}
}

func handleBookError(c *gin.Context, err error) {
	switch err {
	case entity.ErrBookNotFound:
		c.JSON(http.StatusNotFound, generated.ErrorResponse{
			Error: strPtr("book not found"),
			Code:  strPtr("NOT_FOUND"),
		})
	case entity.ErrInvalidBookTitle, entity.ErrInvalidBookAuthor, entity.ErrInvalidBookISBN, entity.ErrInvalidTotalCopies:
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr(err.Error()),
			Code:  strPtr("VALIDATION_ERROR"),
		})
	default:
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Error: strPtr("internal server error"),
			Code:  strPtr("INTERNAL_ERROR"),
		})
	}
}

func handleLoanError(c *gin.Context, err error) {
	switch err {
	case entity.ErrLoanNotFound:
		c.JSON(http.StatusNotFound, generated.ErrorResponse{
			Error: strPtr("loan not found"),
			Code:  strPtr("NOT_FOUND"),
		})
	case entity.ErrUserNotFound:
		c.JSON(http.StatusNotFound, generated.ErrorResponse{
			Error: strPtr("user not found"),
			Code:  strPtr("NOT_FOUND"),
		})
	case entity.ErrBookNotFound:
		c.JSON(http.StatusNotFound, generated.ErrorResponse{
			Error: strPtr("book not found"),
			Code:  strPtr("NOT_FOUND"),
		})
	case entity.ErrBookNotAvailable:
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("book is not available - all copies are borrowed"),
			Code:  strPtr("BOOK_UNAVAILABLE"),
		})
	case entity.ErrUserDisabled:
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("user is disabled"),
			Code:  strPtr("USER_DISABLED"),
		})
	case entity.ErrLoanAlreadyReturned:
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("loan has already been returned"),
			Code:  strPtr("ALREADY_RETURNED"),
		})
	case entity.ErrUserHasActiveLoan:
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("user already has an active loan for this book"),
			Code:  strPtr("ACTIVE_LOAN_EXISTS"),
		})
	default:
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Error: strPtr("internal server error"),
			Code:  strPtr("INTERNAL_ERROR"),
		})
	}
}
