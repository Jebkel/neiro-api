package services

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"neiro-api/internal/helpers"
	"neiro-api/internal/services/i18n"
	"neiro-api/internal/services/jwt"
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

type Service struct {
	Translator
	JwtManager
}

func NewService() *Service {
	return &Service{
		Translator: i18n.NewI18NService(),
		JwtManager: jwt.NewJwtService(),
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
