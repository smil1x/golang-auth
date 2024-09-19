package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"golang-auth/internal/app/dto"
	"golang-auth/internal/app/repository"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
	"net/smtp"
	"strings"
)

type AuthService struct {
	repo repository.Authorization
	jwt  JWT
}

func NewAuthService(repo *repository.Repository, jwt *JWTService) *AuthService {
	return &AuthService{
		repo: repo,
		jwt:  jwt,
	}
}

func (s *AuthService) GetTokens(userId string, userIp string) (*dto.TokensDTO, error) {
	_, err := s.repo.GetUser(userId)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwt.CreateAccessToken(userIp, userId)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwt.CreateRefreshToken(userIp, userId)
	if err != nil {
		return nil, err
	}

	refreshHash, err := s.sha256Hash(refreshToken)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.UpdateRefreshToken(userId, refreshHash)

	return &dto.TokensDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *AuthService) toSha256(str string) [32]byte {
	return sha256.Sum256([]byte(str))
}

func (s *AuthService) sha256Hash(str string) (string, error) {
	strSha256 := s.toSha256(str)

	hash, err := bcrypt.GenerateFromPassword(strSha256[:], bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hash), nil
}

func (s *AuthService) GetIp(r *http.Request) string {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String()
		}
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1"
		}
		return ip
	}

	return ""
}

func (s *AuthService) RefreshTokens(tokens *dto.TokensDTO, userIp string) (*dto.TokensDTO, error) {
	accessClaims, err := s.jwt.ParseAndVerifyAccessToken(tokens.AccessToken)
	if err != nil {
		return nil, err
	}
	refreshClaims, err := s.jwt.ParseAndValidateRefreshToken(tokens.RefreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUser(accessClaims.UserId)
	if err != nil {
		return nil, err
	}

	toSha256 := s.toSha256(tokens.RefreshToken)

	err = bcrypt.CompareHashAndPassword([]byte(user.RefreshHash), toSha256[:])
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if refreshClaims.IP != accessClaims.IP || accessClaims.IP != userIp {
		s.sendWarningEmail(user.Email, refreshClaims.IP, userIp)
	}

	return s.GetTokens(refreshClaims.UserId, userIp)
}

func (s *AuthService) sendWarningEmail(email, oldIP, newIP string) {
	from := "noreply@go-auth.com"
	to := email
	subject := "Warning: IP Address Changed"
	body := fmt.Sprintf("Your IP address has changed from %s to %s. If this wasn't you, please secure your account.", oldIP, newIP)
	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, subject, body)

	smtpHost := "smtp.yourapp.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", "your-email@example.com", "your-password", smtpHost)

	_ = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))

	fmt.Println("Warning email sent to", email)
}
