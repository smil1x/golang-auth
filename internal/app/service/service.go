package service

import (
	"golang-auth/internal/app/dto"
	"golang-auth/internal/app/repository"
	"net/http"
)

type Authorization interface {
	GetTokens(userId string, userIp string) (*dto.TokensDTO, error)
	GetIp(r *http.Request) string
	RefreshTokens(tokens *dto.TokensDTO, userIp string) (*dto.TokensDTO, error)
}

type JWT interface {
	CreateAccessToken(ip string, userId string) (string, error)
	CreateRefreshToken(ip string, userId string) (string, error)
	ParseAndValidateRefreshToken(refreshToken string) (*TokenClaims, error)
	ParseAndVerifyAccessToken(accessToken string) (*TokenClaims, error)
}

type Service struct {
	Authorization
}

func NewService(repository *repository.Repository, config *JWTConfig) *Service {
	return &Service{
		Authorization: NewAuthService(repository, NewJWTService(config)),
	}
}
