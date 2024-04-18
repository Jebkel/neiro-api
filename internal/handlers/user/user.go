package user

import (
	"github.com/gin-gonic/gin"
	"neiro-api/internal/database"
	"neiro-api/internal/handlers"
	"neiro-api/internal/helpers"
	"neiro-api/internal/middlewares"
	"neiro-api/internal/models"
	"net/http"
)

type HandlerUser struct {
	*handlers.Handler
}

func (h HandlerUser) Init(g *gin.RouterGroup) {
	// Initialize handlers
	g.Use(middlewares.JwtAuthMiddleware())
	g.POST("/logout", h.Logout)
}

func (h HandlerUser) Logout(c *gin.Context) {
	db := database.GetDB()

	jwtClaims := helpers.GetJwtClaims(c)

	// Удаление сессии из бд, что бы по токенам нельзя было авторизоаваться
	db.Delete(&models.UserSessions{}, jwtClaims.ID)
	c.AbortWithStatus(http.StatusNoContent)
}
