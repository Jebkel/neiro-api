package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"neiro-api/internal/database"
	"neiro-api/internal/models"
	"neiro-api/internal/services/jwt"
	"net/http"
	"regexp"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := parseJwtFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
		jwtService := jwt.NewJwtService()
		claims, err := jwtService.ParseJwtToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
		if status, err := jwtService.CheckTokenInDB(claims.ID); err != nil || !status {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "session expired",
			})
			return
		}

		var user models.User
		db := database.GetDB()
		if err := db.Find(&user, claims.Subject).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "session expired",
			})
			return
		}
		c.Set("user", user)
		c.Set("language", user.Language)
		c.Set("jwtClaims", claims)
		c.Next()
	}
}

func parseJwtFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	r, err := regexp.Compile("^Bearer (.+)$")
	if err != nil {
		return "", errors.New("invalid token")
	}
	match := r.FindStringSubmatch(authHeader)

	if len(match) != 2 {
		return "", errors.New("authorization header format is invalid")
	}

	tokenString := match[1]
	if len(tokenString) == 0 {
		return "", errors.New("authorization jwt token is empty")
	}

	return tokenString, nil
}
