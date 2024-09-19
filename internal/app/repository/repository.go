package repository

import (
	"database/sql"
	"golang-auth/internal/model"
)

type Authorization interface {
	GetUser(userid string) (*model.User, error)
	UpdateRefreshToken(userId string, refreshToken string) (*string, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
