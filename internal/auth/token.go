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

func (t *TokenService) ExtractUserID(token string) (int, error) {
	jwtToken, err := t.jwtSvc.VerifyToken(token)
	if err != nil || !jwtToken.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id")
	}

	return int(userID), nil
}
