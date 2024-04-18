package auth

import (
	"github.com/gin-gonic/gin"
	"neiro-api/internal/database"
	"neiro-api/internal/handlers"
	"neiro-api/internal/helpers"
	"neiro-api/internal/models"
	"neiro-api/internal/utils"
	"net/http"
)

type HandlerAuth struct {
	*handlers.Handler
}

func (h HandlerAuth) Init(g *gin.RouterGroup) {
	// Initialize handlers
	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
}

func (h HandlerAuth) Register(c *gin.Context) {
	type RequestBody struct {
		Username string `json:"username" binding:"required,gte=6,lte=32"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,gte=8,lte=64"`
	}
	var body RequestBody

	if err := c.ShouldBindJSON(&body); err != nil {
		h.Services.HandleError(c, err)
		return
	}
	var user models.User
	db := database.GetDB()
	if db.Where("username = ?", body.Username).First(&user).Error == nil {
		utils.NewErrorResponse(c, http.StatusBadRequest,
			h.Services.TranslateMessage(helpers.GetLanguage(c), "email_already_using"), &gin.H{})
		return
	}

	user = models.User{
		Email:    body.Email,
		Username: body.Username,
		Language: helpers.GetLanguage(c),
		Password: body.Password,
	}

	err := user.HashPassword()
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Internal server error", &gin.H{})
	}

	// Сохранение пользователя в базе данных
	if err := db.Create(&user).Error; err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Internal server error", &gin.H{})
		return
	}

	// Генерация JWT токенов
	accessToken, refreshToken, err := h.Services.JwtManager.CreateJwtToken(user.ID, c.ClientIP())
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Internal server error", &gin.H{})
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, &gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          &user,
	}, "Register successfully")
}

func (h HandlerAuth) Login(c *gin.Context) {
	type RequestBody struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var body RequestBody

	if err := c.ShouldBindJSON(&body); err != nil {
		h.Services.HandleError(c, err)
		return
	}
	query := "username = ?"
	if helpers.IsEmail(body.Login) {
		query = "email = ?"
	}
	db := database.GetDB()

	var user models.User
	language := helpers.GetLanguage(c)
	if err := db.Where(query, body.Login).First(&user).Error; err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest,
			h.Services.TranslateMessage(language, "invalid_login_or_password"), &gin.H{})
		return
	}

	// Проверка пароля
	check, err := user.ValidatePassword(body.Password)
	if err != nil || !check {
		utils.NewErrorResponse(c, http.StatusBadRequest,
			h.Services.TranslateMessage(language, "invalid_login_or_password"), &gin.H{})
	}

	// Генерация JWT токенов
	accessToken, refreshToken, err := h.Services.JwtManager.CreateJwtToken(user.ID, c.ClientIP())
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Internal server error", &gin.H{})
		return
	}
	utils.NewSuccessResponse(c, http.StatusOK, &gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          &user,
	}, "Login successfully")
}
