package helpers

import (
	"github.com/gin-gonic/gin"
	"neiro-api/internal/models"
	"neiro-api/internal/services/jwt"
	"net/mail"
)

func IsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GetLanguage(c *gin.Context) string {
	lang, ok := c.Get("language")
	if !ok {
		return "en"
	}
	langString, ok := lang.(string)
	if !ok {
		return "en"
	}
	return langString
}

func GetJwtClaims(c *gin.Context) *jwt.CustomClaims {
	claims, ok := c.Get("jwtClaims")
	if !ok {
		return nil
	}
	claimsObject, ok := claims.(*jwt.CustomClaims)
	if !ok {
		return nil
	}
	return claimsObject
}

func GetUser(c *gin.Context) *models.User {
	user, ok := c.Get("user")
	if !ok {
		return nil
	}
	userObject, ok := user.(*models.User)
	if !ok {
		return nil
	}
	return userObject
}
