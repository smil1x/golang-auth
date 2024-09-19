package repository

import (
	"database/sql"
	"errors"
	"golang-auth/internal/model"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) GetUser(userId string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow("SELECT guid, email, refresh_hash FROM users WHERE guid = $1", userId).
		Scan(&u.GUID, &u.Email, &u.RefreshHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return u, nil
}

func (r *AuthPostgres) UpdateRefreshToken(userId string, refreshHash string) (*string, error) {
	u := &model.User{}
	err := r.db.QueryRow("UPDATE users SET refresh_hash = $1 WHERE guid = $2", refreshHash, userId).Scan(&u.GUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u.GUID, nil
}
