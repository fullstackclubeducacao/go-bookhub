package handler

import (
	"log"
	"net/http"

	"bookhub/api/generated"

	"github.com/gin-gonic/gin"
)

// Auth handlers

func (h *Handler) Login(c *gin.Context) {
	var req generated.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Error: strPtr("invalid request body"),
			Code:  strPtr("BAD_REQUEST"),
		})
		return
	}

	user, err := h.userUseCase.ValidateCredentials(c.Request.Context(), string(req.Email), req.Password)
	if err != nil {
		log.Printf("Login failed for %s: %v", req.Email, err)
		c.JSON(http.StatusUnauthorized, generated.ErrorResponse{
			Error: strPtr("invalid credentials"),
			Code:  strPtr("UNAUTHORIZED"),
		})
		return
	}

	token, expiresAt, err := h.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Error: strPtr("failed to generate token"),
			Code:  strPtr("INTERNAL_ERROR"),
		})
		return
	}

	c.JSON(http.StatusOK, generated.LoginResponse{
		Token:     strPtr(token),
		ExpiresAt: timePtr(expiresAt),
		User:      userToResponse(user),
	})
}
