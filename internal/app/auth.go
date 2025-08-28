package app

import (
	"errors"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/repository"
	"github.com/slipe-fun/skid-backend/internal/service"
)

type AuthApp struct {
	users  *repository.UserRepo
	jwtSvc *service.JWTService
}

func NewAuthApp(users *repository.UserRepo, jwt *service.JWTService) *AuthApp {
	return &AuthApp{
		users:  users,
		jwtSvc: jwt,
	}
}

func (a *AuthApp) Login(username, password string, expire time.Duration) (string, *domain.User, error) {
	user, err := a.users.GetByUsername(username)

	if err != nil {
		return "", nil, errors.New("user not found")
	}

	ok, err := service.VerifyPassword(password, user.Password)

	if err != nil || !ok {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := a.jwtSvc.GenerateToken(user.ID, expire)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (a *AuthApp) Register(username, password string, expire time.Duration) (string, *domain.User, error) {
	_, err := a.users.GetByUsername(username)

	if err == nil {
		return "", nil, errors.New("user already exists")
	}

	hashedPassword, err := service.HashPassword(password)

	if err != nil {
		return "", nil, err
	}

	user, err := a.users.Create(&domain.User{
		Username: username,
		Password: hashedPassword,
	})

	if err != nil {
		return "", nil, err
	}

	token, err := a.jwtSvc.GenerateToken(user.ID, expire)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
