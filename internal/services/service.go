package services

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"neiro-api/internal/helpers"
	"neiro-api/internal/redis"
	"neiro-api/internal/services/i18n"
	"neiro-api/internal/services/jwt"
	"neiro-api/internal/services/mail"
	"neiro-api/internal/utils"
	"net/http"
)

type Translator interface {
	TranslateValidationError(language string, err validator.ValidationErrors) *gin.H
	TranslateMessage(language string, key string) string
}

type JwtManager interface {
	ParseJwtToken(tokenString string) (*jwt.CustomClaims, error)
	CheckTokenInDB(jwtID string) (bool, error)
	DeprecateSession(jwtID string)
	CreateJwtToken(userID uint, ipAddress string) (signedToken string, refreshToken string, err error)
}

type MailManager interface {
	New() *mail.ServiceMail
	To(recipient string) *mail.ServiceMail
	From(from string) *mail.ServiceMail
	Line(line string) *mail.ServiceMail
	Header(header string) *mail.ServiceMail
	Subject(subject string) *mail.ServiceMail
	Send()
}

type Service struct {
	Translator
	JwtManager
	MailManager
	RedisManager *redis.ManagerRedis
}

func NewService() *Service {
	return &Service{
		Translator:   i18n.NewI18NService(),
		JwtManager:   jwt.NewJwtService(),
		MailManager:  mail.NewMailService(),
		RedisManager: redis.GetRedis(),
	}
}

func (s *Service) HandleError(c *gin.Context, err error) {
	language := helpers.GetLanguage(c)
	err = unwrapRecursive(err)
	switch errs := err.(type) {
	case validator.ValidationErrors:
		utils.NewErrorResponse(c, http.StatusBadRequest, "validation error", s.TranslateValidationError(language, errs))
	case *json.SyntaxError:
		utils.NewErrorResponse(c, http.StatusBadRequest, "bad json syntax", &gin.H{})
	default:
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Internal server error", &gin.H{})
	}
}

func unwrapRecursive(err error) error {
	for {
		internalErr := errors.Unwrap(err)
		if internalErr == nil {
			break
		}
		err = internalErr
	}
	return err
}
