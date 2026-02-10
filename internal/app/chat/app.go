package chat

import (
	"github.com/slipe-fun/skid-backend/internal/app/session"
	"github.com/slipe-fun/skid-backend/internal/auth"
	"github.com/slipe-fun/skid-backend/internal/repository/chat"
)

type ChatApp struct {
	sessionApp *session.SessionApp
	chats      *chat.ChatRepo
	tokenSvc   *auth.TokenService
}

func NewChatApp(sessionApp *session.SessionApp, chats *chat.ChatRepo, tokenSvc *auth.TokenService) *ChatApp {
	return &ChatApp{
		sessionApp: sessionApp,
		chats:      chats,
		tokenSvc:   tokenSvc,
	}
}
