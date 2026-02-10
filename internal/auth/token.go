package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	jwtSvc *JWTService
}

func NewTokenService(jwtSvc *JWTService) *TokenService {
	return &TokenService{jwtSvc: jwtSvc}
}

func (t *TokenService) ExtractUserID(tokenStr string) (int, error) {
	token, err := t.jwtSvc.VerifyToken(tokenStr)
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id")
	}

	return int(userID), nil
}
