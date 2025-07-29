package config

import (
	"database/sql"
)

type AppConfig struct {
	DB *sql.DB
}

func GetAppConfig(db *sql.DB) *AppConfig {
	return &AppConfig{DB: db}
}
