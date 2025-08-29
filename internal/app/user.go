package app

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	UserRepo "github.com/slipe-fun/skid-backend/internal/repository/user"
	"github.com/slipe-fun/skid-backend/internal/service"
)

type UserApp struct {
	users    *UserRepo.UserRepo
	jwtSvc   *service.JWTService
	tokenSvc *service.TokenService
}

func NewUserApp(users *UserRepo.UserRepo, jwtSvc *service.JWTService, tokenSvc *service.TokenService) *UserApp {
	return &UserApp{
		users:    users,
		jwtSvc:   jwtSvc,
		tokenSvc: tokenSvc,
	}
}

func (u *UserApp) GetUserByToken(tokenStr string) (*domain.User, error) {
	userID, err := u.tokenSvc.ExtractUserID(tokenStr)
	if err != nil {
		return nil, err
	}

	user, err := u.users.GetById(int(userID))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserApp) GetUserById(id int) (*domain.User, error) {
	user, err := u.users.GetById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
