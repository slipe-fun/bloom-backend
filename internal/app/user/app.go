package user

import (
	"github.com/slipe-fun/skid-backend/internal/app/session"
	"github.com/slipe-fun/skid-backend/internal/auth"
	"github.com/slipe-fun/skid-backend/internal/repository/user"
)

type UserApp struct {
	sessionApp *session.SessionApp
	users      *user.UserRepo
	jwtSvc     *auth.JWTService
	tokenSvc   *auth.TokenService
}

func NewUserApp(sessionApp *session.SessionApp,
	users *user.UserRepo,
	jwtSvc *auth.JWTService,
	tokenSvc *auth.TokenService) *UserApp {
	return &UserApp{
		sessionApp: sessionApp,
		users:      users,
		jwtSvc:     jwtSvc,
		tokenSvc:   tokenSvc,
	}
}
