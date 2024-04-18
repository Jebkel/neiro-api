package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"neiro-api/config"
	"neiro-api/internal/database"
	"neiro-api/internal/models"
	"neiro-api/internal/redis"
	"neiro-api/internal/routes"
	"neiro-api/internal/utils"
)

func main() {
	err := config.Init("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	utils.UseJSONLogFormat()

	sqlDB := database.Init()
	models.Migrate()

	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Fatalf("failed to close database connection: %v", err)
		}
	}(sqlDB)

	redis.Init()

	gin.SetMode(gin.DebugMode)

	r := routes.Init()
	err = r.SetTrustedProxies(nil)
	if err != nil {
		log.Error(err)
	}
	cfg := config.GetConfig()
	log.Fatal(r.Run(fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)))
}
