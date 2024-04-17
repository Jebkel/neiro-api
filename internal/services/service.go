package services

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"neiro-api/internal/utils"
	"net/http"
)

type Translator interface {
	TranslateValidationError(language string, err validator.ValidationErrors) *gin.H
}

type Service struct {
	Translator
}

func NewService() *Service {
	return &Service{
		Translator: NewI18NService(),
	}
}

func (s *Service) HandleError(c *gin.Context, err error) {
	language := getLanguage(c)
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

func getLanguage(c *gin.Context) string {
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
