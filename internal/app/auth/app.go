package auth

import (
	"github.com/slipe-fun/skid-backend/internal/app/session"
	verificationapp "github.com/slipe-fun/skid-backend/internal/app/verification"
	"github.com/slipe-fun/skid-backend/internal/auth"
	"github.com/slipe-fun/skid-backend/internal/oauth/google"
	"github.com/slipe-fun/skid-backend/internal/repository/user"
	verificationrepo "github.com/slipe-fun/skid-backend/internal/repository/verification"
)

type AuthApp struct {
	sessionApp *session.SessionApp
	users      *user.UserRepo
	codesRepo  *verificationrepo.VerificationRepo
	codesApp   *verificationapp.VerificationApp
	jwtSvc     *auth.JWTService
	google     *google.GoogleAuthService
}

func NewAuthApp(sessionApp *session.SessionApp,
	users *user.UserRepo,
	codesRepo *verificationrepo.VerificationRepo,
	codesApp *verificationapp.VerificationApp,
	jwt *auth.JWTService,
	google *google.GoogleAuthService) *AuthApp {
	return &AuthApp{
		sessionApp: sessionApp,
		users:      users,
		codesRepo:  codesRepo,
		codesApp:   codesApp,
		jwtSvc:     jwt,
		google:     google,
	}
}
