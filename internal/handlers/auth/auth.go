package auth

import (
	"github.com/gin-gonic/gin"
	"neiro-api/internal/handlers"
	"net/http"
)

type HandlerAuth struct {
	*handlers.Handler
}

func (h HandlerAuth) Init(g *gin.RouterGroup) {
	// Initialize handlers
	g.POST("/register", h.Register)
}

func (h HandlerAuth) Register(c *gin.Context) {
	type RequestBody struct {
		Username string `json:"username" binding:"required,gte=6,lte=32"`
		Password string `json:"password" binding:"required,gte=8,lte=64"`
	}
	var body RequestBody

	if err := c.ShouldBindJSON(&body); err != nil {
		h.Services.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
