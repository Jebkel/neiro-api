package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"neiro-api/config"
	"neiro-api/internal/database"
	"neiro-api/internal/models"
	"strconv"
	"time"
)

type JwtService struct {
	db          *gorm.DB
	secretKey   []byte
	tokenExpire time.Duration
}

type JwtCustomClaims struct {
	jwt.RegisteredClaims
}

func NewJwtService() *JwtService {
	cfg := config.GetConfig().JwtConfig
	return &JwtService{
		db:          database.GetDB(),
		secretKey:   []byte(cfg.JWTSecret),
		tokenExpire: cfg.JwtDuration,
	}
}

// ParseJwtToken : Проверяет рабочий ли токен и возвращает данные из токена в случае успеха
func (j *JwtService) ParseJwtToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка на метод, используемый для подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signed method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok || claims == nil {
		log.Error("failed to parse JWT claims")
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// CheckTokenInDB : проверяет, существует ли сессия в базе данных
func (j *JwtService) CheckTokenInDB(jwtID string) (bool, error) {
	var UserSession models.UserSessions
	result := j.db.Where("id = ?", jwtID).First(&UserSession)
	if result.Error != nil || result.RowsAffected == 0 {
		return false, errors.New("session expired")
	}
	return true, nil
}

// DeprecateSession : удаляет сессию из базы данных
func (j *JwtService) DeprecateSession(jwtID string) {
	j.db.Delete(&models.UserSessions{}, jwtID)
}

func (j *JwtService) CreateJwtToken(userID uint, ipAddress string) (signedToken string, refreshToken string, err error) {
	refreshToken = uuid.New().String()
	userSession := models.UserSessions{
		IpAddress:    ipAddress,
		UserID:       userID,
		RefreshToken: refreshToken,
	}
	j.db.Create(&userSession)

	claims := JwtCustomClaims{
		jwt.RegisteredClaims{
			ID:        strconv.Itoa(int(userSession.ID)),
			Subject:   strconv.Itoa(int(userID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExpire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString(j.secretKey)
	if err != nil {
		return "", "", err
	}
	return signedToken, refreshToken, err
}
