package middleware

import (
	"net/http"
	"strings"

	"bookhub/api/generated"
	"bookhub/internal/infrastructure/auth"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "user_id"
	UserEmailKey        = "user_email"
)

// JWTAuthWithOpenAPI creates a middleware that checks authentication based on OpenAPI spec.
// It only validates JWT tokens for routes that have security defined in the OpenAPI specification.
// Routes without security (like /auth/login) will pass through without authentication.
func JWTAuthWithOpenAPI(jwtService auth.JWTService) generated.MiddlewareFunc {
	return func(c *gin.Context) {
		// Check if this route requires authentication by looking for BearerAuthScopes
		// The oapi-codegen sets this value only for routes with security defined in OpenAPI
		_, exists := c.Get(generated.BearerAuthScopes)
		if !exists {
			// Route does not require authentication
			c.Next()
			return
		}

		// Route requires authentication - validate JWT token
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
				"code":  "UNAUTHORIZED",
			})
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
				"code":  "UNAUTHORIZED",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
				"code":  "UNAUTHORIZED",
			})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Next()
	}
}
