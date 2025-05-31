package db

import (
	"fmt"
	"go_link_reducer/config"
	"go_link_reducer/types"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	var dsn string

	if cfg.AppEnvironment == "development" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			cfg.PublicHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.Port,
		)
		log.Print("using development dsn")
	} else {
		dsn = cfg.DatabaseURL
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	models := []interface{}{
		&types.URL{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return nil, err
	}

	return db, nil
}
