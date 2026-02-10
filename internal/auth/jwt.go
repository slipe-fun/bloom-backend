package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/slipe-fun/skid-backend/internal/config"
)

type JWTService struct {
	secret string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: secret}
}

func (s *JWTService) GenerateToken(userID int) (string, error) {
	cfg := config.LoadConfig("configs/config.yaml")
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(cfg.JWTExpireDuration()).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTService) VerifyToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenInvalidSubject
		}
		return []byte(s.secret), nil
	})
}
