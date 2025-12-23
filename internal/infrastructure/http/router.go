package http

import (
	"encoding/json"
	"log"
	"net/http"

	"bookhub/api/generated"
	"bookhub/internal/infrastructure/http/handler"
	"bookhub/internal/infrastructure/http/middleware"

	"github.com/flowchartsman/swaggerui"
	"github.com/gin-gonic/gin"
)

func NewRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	swagger, err := generated.GetSwagger()
	if err != nil {
		log.Fatalf("failed to get swagger spec: %v", err)
	}
	swaggerJSON, err := json.Marshal(swagger)
	if err != nil {
		log.Fatalf("failed to marshal swagger spec: %v", err)
	}
	router.GET("/docs/*any", gin.WrapH(http.StripPrefix("/docs", swaggerui.Handler(swaggerJSON))))

	api := router.Group("/api/v1")
	{
		// Register all OpenAPI-generated handlers with JWT authentication middleware
		// The middleware checks BearerAuthScopes from OpenAPI spec to determine if auth is required
		generated.RegisterHandlersWithOptions(api, h, generated.GinServerOptions{
			Middlewares: []generated.MiddlewareFunc{
				middleware.JWTAuthWithOpenAPI(h.JWTService()),
			},
		})
	}

	return router
}
