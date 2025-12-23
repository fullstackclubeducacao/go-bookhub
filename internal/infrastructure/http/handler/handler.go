package handler

import (
	"bookhub/internal/infrastructure/auth"
	"bookhub/internal/usecase"
)

type Handler struct {
	userUseCase usecase.UserUseCase
	bookUseCase usecase.BookUseCase
	loanUseCase usecase.LoanUseCase
	jwtService  auth.JWTService
}

func NewHandler(
	userUseCase usecase.UserUseCase,
	bookUseCase usecase.BookUseCase,
	loanUseCase usecase.LoanUseCase,
	jwtService auth.JWTService,
) *Handler {
	return &Handler{
		userUseCase: userUseCase,
		bookUseCase: bookUseCase,
		loanUseCase: loanUseCase,
		jwtService:  jwtService,
	}
}

func (h *Handler) JWTService() auth.JWTService {
	return h.jwtService
}
