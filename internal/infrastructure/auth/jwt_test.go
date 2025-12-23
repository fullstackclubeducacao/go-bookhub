package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTService_GenerateToken(t *testing.T) {
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: 1 * time.Hour,
		Issuer:        "bookhub-test",
	}

	service := NewJWTService(config)

	t.Run("generate valid token", func(t *testing.T) {
		userID := uuid.New()
		email := "test@example.com"

		token, expiresAt, err := service.GenerateToken(userID, email)
		if err != nil {
			t.Errorf("JWTService.GenerateToken() unexpected error = %v", err)
			return
		}

		if token == "" {
			t.Error("JWTService.GenerateToken() token should not be empty")
		}

		if expiresAt.Before(time.Now()) {
			t.Error("JWTService.GenerateToken() expiresAt should be in the future")
		}
	})
}

func TestJWTService_ValidateToken(t *testing.T) {
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: 1 * time.Hour,
		Issuer:        "bookhub-test",
	}

	service := NewJWTService(config)

	t.Run("validate valid token", func(t *testing.T) {
		userID := uuid.New()
		email := "test@example.com"

		token, _, _ := service.GenerateToken(userID, email)

		claims, err := service.ValidateToken(token)
		if err != nil {
			t.Errorf("JWTService.ValidateToken() unexpected error = %v", err)
			return
		}

		if claims.UserID != userID {
			t.Errorf("JWTService.ValidateToken() userID = %v, want %v", claims.UserID, userID)
		}
		if claims.Email != email {
			t.Errorf("JWTService.ValidateToken() email = %v, want %v", claims.Email, email)
		}
	})

	t.Run("validate invalid token", func(t *testing.T) {
		_, err := service.ValidateToken("invalid-token")
		if err != ErrInvalidToken {
			t.Errorf("JWTService.ValidateToken() error = %v, wantErr %v", err, ErrInvalidToken)
		}
	})

	t.Run("validate expired token", func(t *testing.T) {
		expiredConfig := JWTConfig{
			SecretKey:     "test-secret-key",
			TokenDuration: -1 * time.Hour,
			Issuer:        "bookhub-test",
		}
		expiredService := NewJWTService(expiredConfig)

		userID := uuid.New()
		email := "test@example.com"
		token, _, _ := expiredService.GenerateToken(userID, email)

		_, err := service.ValidateToken(token)
		if err != ErrExpiredToken {
			t.Errorf("JWTService.ValidateToken() error = %v, wantErr %v", err, ErrExpiredToken)
		}
	})

	t.Run("validate token with wrong secret", func(t *testing.T) {
		userID := uuid.New()
		email := "test@example.com"
		token, _, _ := service.GenerateToken(userID, email)

		wrongConfig := JWTConfig{
			SecretKey:     "wrong-secret-key",
			TokenDuration: 1 * time.Hour,
			Issuer:        "bookhub-test",
		}
		wrongService := NewJWTService(wrongConfig)

		_, err := wrongService.ValidateToken(token)
		if err != ErrInvalidToken {
			t.Errorf("JWTService.ValidateToken() error = %v, wantErr %v", err, ErrInvalidToken)
		}
	})
}
