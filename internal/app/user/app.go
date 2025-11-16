package UserApp

import (
	SessionApp "github.com/slipe-fun/skid-backend/internal/app/session"
	UserRepo "github.com/slipe-fun/skid-backend/internal/repository/user"
	"github.com/slipe-fun/skid-backend/internal/service"
)

type UserApp struct {
	sessionApp *SessionApp.SessionApp
	users      *UserRepo.UserRepo
	jwtSvc     *service.JWTService
	tokenSvc   *service.TokenService
}

func NewUserApp(sessionApp *SessionApp.SessionApp,
	users *UserRepo.UserRepo,
	jwtSvc *service.JWTService,
	tokenSvc *service.TokenService) *UserApp {
	return &UserApp{
		sessionApp: sessionApp,
		users:      users,
		jwtSvc:     jwtSvc,
		tokenSvc:   tokenSvc,
	}
}
