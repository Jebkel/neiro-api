package passwordReset

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"neiro-api/internal/database"
	"neiro-api/internal/handlers"
	"neiro-api/internal/models"
	redis2 "neiro-api/internal/redis"
	"time"
)

type HandlerPasswordReset struct {
	*handlers.Handler
}

func (h HandlerPasswordReset) Init(g *gin.RouterGroup) {
	// Initialize handlers
	g.POST("/code/send", h.SendResetCode)
	g.POST("/code/validate", h.ValidateResetCode)
}

func (h HandlerPasswordReset) SendResetCode(c *gin.Context) {
	type RequestBody struct {
		Email string `json:"email" binding:"required,email"`
	}
	var body RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		h.Services.HandleError(c, err)
		return
	}
	db := database.GetDB()
	var user models.User

	if result := db.Where("email = ?", body.Email).First(&user); result.Error != nil || result.RowsAffected == 0 {
		// Отправляем статус 204, что бы нельзя было перебирать существующих пользователей
		log.Debug(result.Error, result.RowsAffected)
		c.AbortWithStatus(204)
		return
	}

	code := generateRandomCode()

	h.Services.MailManager.New().To(user.Email).Subject("Востановление пароля").
		Header("Востановление пароля").Line("Кто-то запросил востонавление пароля").
		Line(fmt.Sprintf("Ваш код для востонавления: %s", code)).
		Line("Если это были не вы, смело игнорируйте это сообщение").Send()

	redisManager := h.Services.RedisManager

	redisManager.ClientRedis.Set(redisManager.Ctx, fmt.Sprintf("resetCodes:email_%s", body.Email), code, time.Minute*60)

	c.AbortWithStatus(204)
}

func (h HandlerPasswordReset) ValidateResetCode(c *gin.Context) {
	type RequestBody struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}
	var body RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		h.Services.HandleError(c, err)
		return
	}
	db := database.GetDB()
	var user models.User

	if result := db.Where("email = ?", body.Email).First(&user); result.Error != nil || result.RowsAffected == 0 {
		// Отправляем статус 204, что бы нельзя было перебирать существующих пользователей
		log.Debug(result.Error, result.RowsAffected)
		c.AbortWithStatus(400)
		return
	}
	redisManager := h.Services.RedisManager

	if validateCode(redisManager, body.Code, body.Email) {
		c.AbortWithStatus(204)
	} else {
		c.AbortWithStatus(400)
	}
}

func (h HandlerPasswordReset) ResetPassword(c *gin.Context) {
	type RequestBody struct {
		Email       string `json:"email" binding:"required,email"`
		Code        string `json:"code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	var body RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		h.Services.HandleError(c, err)
		return
	}
	db := database.GetDB()
	var user models.User

	if result := db.Where("email = ?", body.Email).First(&user); result.Error != nil || result.RowsAffected == 0 {
		// Отправляем статус 204, что бы нельзя было перебирать существующих пользователей
		log.Debug(result.Error, result.RowsAffected)
		c.AbortWithStatus(400)
		return
	}
	redisManager := h.Services.RedisManager

	if !validateCode(redisManager, body.Code, body.Email) {
		c.AbortWithStatus(400)
	}
	user.Password = body.NewPassword
	err := user.HashPassword()
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	db.Save(user)
	c.AbortWithStatus(204)
}

func validateCode(redisManager *redis2.ManagerRedis, code string, email string) bool {
	rcode, err := redisManager.ClientRedis.Get(redisManager.Ctx,
		fmt.Sprintf("resetCodes:email_%s", email)).Result()
	if errors.Is(err, redis.Nil) {
		return false
	} else if err != nil {
		return false
	}
	if rcode == code {
		return true
	}
	return false
}

func generateRandomCode() string {
	const charset = "0123456789"
	code := make([]byte, 8)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range code {
		code[i] = charset[r.Intn(len(charset))]
	}
	return string(code)
}
