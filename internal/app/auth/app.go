package AuthApp

import (
	VerificationApp "github.com/slipe-fun/skid-backend/internal/app/verification"
	UserRepo "github.com/slipe-fun/skid-backend/internal/repository/user"
	VerificationRepo "github.com/slipe-fun/skid-backend/internal/repository/verification"
	"github.com/slipe-fun/skid-backend/internal/service"
)

type AuthApp struct {
	users     *UserRepo.UserRepo
	codesRepo *VerificationRepo.VerificationRepo
	codesApp  *VerificationApp.VerificationApp
	jwtSvc    *service.JWTService
}

func NewAuthApp(users *UserRepo.UserRepo,
	codesRepo *VerificationRepo.VerificationRepo,
	codesApp *VerificationApp.VerificationApp,
	jwt *service.JWTService) *AuthApp {
	return &AuthApp{
		users:     users,
		codesRepo: codesRepo,
		codesApp:  codesApp,
		jwtSvc:    jwt,
	}
}
