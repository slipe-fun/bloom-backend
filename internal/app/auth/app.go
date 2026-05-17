package auth

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/redis/go-redis/v9"
)

type AuthApp struct {
	sessionApp  SessionApp
	users       UserRepo
	credentials CredentialRepo
	rdb         *redis.Client
	webauthn    *webauthn.WebAuthn
}

func NewAuthApp(
	sessionApp SessionApp,
	users UserRepo,
	credentials CredentialRepo,
	rdb *redis.Client,
	webauthn *webauthn.WebAuthn,
) *AuthApp {
	return &AuthApp{
		sessionApp:  sessionApp,
		users:       users,
		credentials: credentials,
		rdb:         rdb,
		webauthn:    webauthn,
	}
}
