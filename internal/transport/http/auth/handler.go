package auth

import (
	AuthApp "github.com/slipe-fun/skid-backend/internal/app/auth"
	"github.com/slipe-fun/skid-backend/internal/oauth/google"
)

type AuthHandler struct {
	authApp *AuthApp.AuthApp
	google  *google.GoogleAuthService
}

func NewAuthHandler(authApp *AuthApp.AuthApp, google *google.GoogleAuthService) *AuthHandler {
	return &AuthHandler{authApp: authApp, google: google}
}
