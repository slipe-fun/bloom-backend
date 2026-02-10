package session

import (
	"github.com/slipe-fun/skid-backend/internal/auth"
	"github.com/slipe-fun/skid-backend/internal/repository/session"
	"github.com/slipe-fun/skid-backend/internal/repository/user"
)

type SessionApp struct {
	session  *session.SessionRepo
	users    *user.UserRepo
	jwtSvc   *auth.JWTService
	tokenSvc *auth.TokenService
}

func NewSessionApp(session *session.SessionRepo,
	users *user.UserRepo,
	jwtSvc *auth.JWTService,
	tokenSvc *auth.TokenService) *SessionApp {
	return &SessionApp{
		session:  session,
		users:    users,
		jwtSvc:   jwtSvc,
		tokenSvc: tokenSvc,
	}
}
