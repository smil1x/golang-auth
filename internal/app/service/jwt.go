package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTConfig struct {
	PrivateRsaKey      string        `json:"private_rsa_key" validate:"required"`
	PublicRsaKey       string        `json:"public_rsa_key" validate:"required"`
	RefreshTokenSecret string        `json:"refresh_token_secret" validate:"required,min=20,max=32"`
	AccessTokenTTL     time.Duration `json:"access_token_ttl" validate:"required,max=10m"`
	RefreshTokenTTL    time.Duration `json:"refresh_token_ttl" validate:"required,min=30m"`
}

type JWTService struct {
	config *JWTConfig
}

func NewJWTService(cfg *JWTConfig) *JWTService {
	return &JWTService{
		config: cfg,
	}
}

type TokenClaims struct {
	IP     string `json:"ip"`
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *JWTService) CreateAccessToken(ip string, userId string) (string, error) {
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		IP:     ip,
		UserId: userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	signKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(s.config.PrivateRsaKey))
	return token.SignedString(signKey)
}

func (s *JWTService) CreateRefreshToken(ip string, userId string) (string, error) {
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		IP:     ip,
		UserId: userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.RefreshTokenSecret))
}

func (s *JWTService) ParseAndValidateRefreshToken(refreshToken string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, _ := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.RefreshTokenSecret), nil
	})

	if token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (s *JWTService) ParseAndVerifyAccessToken(accessToken string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(s.config.PublicRsaKey))

	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})

	switch {
	case token.Valid:
		return claims, nil
	case errors.Is(err, jwt.ErrTokenMalformed):

		return nil, errors.New("access token is malformed")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, errors.New("invalid access token signature")
	case errors.Is(err, jwt.ErrTokenExpired):
		return claims, nil
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, errors.New("token is not valid yet")
	default:
		return nil, errors.New("invalid token")
	}

}
