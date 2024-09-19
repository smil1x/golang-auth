package repository

import (
	"database/sql"
	"fmt"
)

type Config struct {
	Host     string `json:"host" validate:"required,hostname"`
	Port     string `json:"port" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	DBName   string `json:"dbName" validate:"required"`
	SSLMode  string `json:"sslMode" validate:"required"`
}

func ConnectPostgresDB(cfg Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}
