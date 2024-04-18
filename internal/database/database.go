package database

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"neiro-api/config"
	"time"
)

var db *gorm.DB

func Init() *sql.DB {
	var err error
	cfg := config.GetConfig().DBConfig
	switch cfg.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Collation)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open("internal/database/database.sqlite"), &gorm.Config{})
	default:
		panic("unsupported database type")
	}
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(cfg.MinConn)
	sqlDB.SetMaxOpenConns(cfg.MaxConn)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return sqlDB
}

func GetDB() *gorm.DB {
	return db
}
