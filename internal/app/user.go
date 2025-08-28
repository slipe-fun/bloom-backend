package app

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/repository"
	"github.com/slipe-fun/skid-backend/internal/service"
)

type UserApp struct {
	users  *repository.UserRepo
	jwtSvc *service.JWTService
}

func NewUserApp(users *repository.UserRepo, jwtSvc *service.JWTService) *UserApp {
	return &UserApp{users: users, jwtSvc: jwtSvc}
}

func (u *UserApp) GetUserByToken(tokenStr string) (*domain.User, error) {
	token, err := u.jwtSvc.VerifyToken(tokenStr)
	if err != nil || !token.Valid {
		fmt.Print(err)
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id")
	}

	user, err := u.users.GetById(int(userID))
	if err != nil {
		return nil, err
	}

	return user, nil
}
