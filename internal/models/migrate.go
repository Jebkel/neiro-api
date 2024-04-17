package models

import (
	log "github.com/sirupsen/logrus"
	"neiro-api/internal/database"
)

func Migrate() {
	db := database.GetDB()
	err := db.AutoMigrate(&User{}, &UserSessions{})
	if err != nil {
		log.Error(err)
	}
}
