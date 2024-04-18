package utils

import (
	"github.com/gin-gonic/gin"
	"neiro-api/internal/services"
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

func GetJwtClaims(c *gin.Context) *services.JwtCustomClaims {
	claims, ok := c.Get("jwtClaims")
	if !ok {
		return nil
	}
	claimsObject, ok := claims.(*services.JwtCustomClaims)
	if !ok {
		return nil
	}
	return claimsObject
}
