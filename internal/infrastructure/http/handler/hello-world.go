package handler

import (
	"bookhub/api/generated"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) MyHelloWorld(c *gin.Context) {
	result := generated.HelloWorldResponse{
		Title: "Minha rota customizada",
	}
	c.JSON(http.StatusOK, result)
}
