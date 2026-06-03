package auth

import (
	"github.com/redis/go-redis/v9"
)

type AuthApp struct {
	sessionApp SessionApp
	keysApp    KeysApp
	users      UserRepo
	rdb        *redis.Client
}

func NewAuthApp(
	sessionApp SessionApp,
	keysApp KeysApp,
	users UserRepo,
	rdb *redis.Client,
) *AuthApp {
	return &AuthApp{
		sessionApp: sessionApp,
		keysApp:    keysApp,
		users:      users,
		rdb:        rdb,
	}
}
