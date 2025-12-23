package handler

import (
	"net/http"

	"bookhub/api/generated"
	"bookhub/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// User handlers

func (h *Handler) ListUsers(c *gin.Context, params generated.ListUsersParams) {
	page := 1
	limit := 10

	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	users, total, err := h.userUseCase.List(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Error: strPtr("failed to list users"),
			Code:  strPtr("INTERNAL_ERROR"),
		})
		return
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, generated.UserListResponse{
		Data:       usersToResponse(users),
		Pagination: paginationResponse(page, limit, total, totalPages),
	})
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req generated.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("invalid request body"),
			Code:  strPtr("BAD_REQUEST"),
		})
		return
	}

	user, err := h.userUseCase.Create(c.Request.Context(), usecase.CreateUserInput{
		Name:     req.Name,
		Email:    string(req.Email),
		Password: req.Password,
	})
	if err != nil {
		handleUserError(c, err)
		return
	}

	c.JSON(http.StatusCreated, generated.UserResponse{
		Data: userToResponse(user),
	})
}

func (h *Handler) GetUserById(c *gin.Context, id openapi_types.UUID) {
	userID, err := uuid.Parse(id.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("invalid user ID"),
			Code:  strPtr("BAD_REQUEST"),
		})
		return
	}

	user, err := h.userUseCase.GetByID(c.Request.Context(), userID)
	if err != nil {
		handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, generated.UserResponse{
		Data: userToResponse(user),
	})
}

func (h *Handler) UpdateUser(c *gin.Context, id openapi_types.UUID) {
	userID, err := uuid.Parse(id.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("invalid user ID"),
			Code:  strPtr("BAD_REQUEST"),
		})
		return
	}

	var req generated.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("invalid request body"),
			Code:  strPtr("BAD_REQUEST"),
		})
		return
	}

	input := usecase.UpdateUserInput{}
	if req.Name != nil {
		input.Name = req.Name
	}
	if req.Email != nil {
		emailStr := string(*req.Email)
		input.Email = &emailStr
	}

	user, err := h.userUseCase.Update(c.Request.Context(), userID, input)
	if err != nil {
		handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, generated.UserResponse{
		Data: userToResponse(user),
	})
}

func (h *Handler) DisableUser(c *gin.Context, id openapi_types.UUID) {
	userID, err := uuid.Parse(id.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("invalid user ID"),
			Code:  strPtr("BAD_REQUEST"),
		})
		return
	}

	if err := h.userUseCase.Disable(c.Request.Context(), userID); err != nil {
		handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, generated.MessageResponse{
		Message: strPtr("user disabled successfully"),
	})
}
