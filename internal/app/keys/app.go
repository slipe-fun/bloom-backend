package keys

import (
	"github.com/slipe-fun/skid-backend/internal/app/chat"
	"github.com/slipe-fun/skid-backend/internal/app/session"
	"github.com/slipe-fun/skid-backend/internal/app/user"
	"github.com/slipe-fun/skid-backend/internal/repository/keys"
)

type KeysApp struct {
	sessionApp *session.SessionApp
	keys       *keys.KeysRepo
	users      *user.UserApp
	chats      *chat.ChatApp
}

func NewKeysApp(sessionApp *session.SessionApp,
	keys *keys.KeysRepo,
	users *user.UserApp,
	chats *chat.ChatApp) *KeysApp {
	return &KeysApp{
		sessionApp: sessionApp,
		keys:       keys,
		users:      users,
		chats:      chats,
	}
}
