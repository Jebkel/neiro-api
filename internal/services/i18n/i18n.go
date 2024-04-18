package i18n

import (
	"fmt"
	"github.com/eduardolat/goeasyi18n"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gobeam/stringy"
	log "github.com/sirupsen/logrus"
	"os"
)

type I18nService struct {
	I18n *goeasyi18n.I18n
}

func NewI18NService() *I18nService {
	i18n := goeasyi18n.NewI18n()

	entries, err := os.ReadDir("locales/")
	if err != nil {
		panic("no found directories on locales dir")
	}

	for _, entry := range entries {
		if entry.IsDir() {
			langData, err := goeasyi18n.LoadFromJsonFiles(fmt.Sprintf("locales/%s/*.json", entry.Name()))
			if err != nil {
				log.Error(err)
			}
			i18n.AddLanguage(entry.Name(), langData)
		}
	}
	return &I18nService{
		I18n: i18n,
	}
}

func (i *I18nService) TranslateValidationError(language string, err validator.ValidationErrors) *gin.H {
	errorMessages := gin.H{}
	for _, e := range err {
		errorMessages[stringy.New(e.Field()).SnakeCase().ToLower()] = i.I18n.T(language,
			fmt.Sprintf("field_%s", e.Tag()),
			goeasyi18n.Options{
				Data: map[string]string{
					"valData": e.Param(),
				},
			})
	}
	return &errorMessages
}

func (i *I18nService) TranslateMessage(language string, key string) string {
	return i.I18n.T(language, key, goeasyi18n.Options{})
}
