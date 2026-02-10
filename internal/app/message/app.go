package message

import (
	"github.com/slipe-fun/skid-backend/internal/app/chat"
	"github.com/slipe-fun/skid-backend/internal/app/session"
	"github.com/slipe-fun/skid-backend/internal/auth"
	"github.com/slipe-fun/skid-backend/internal/repository/message"
)

type MessageApp struct {
	sessionApp *session.SessionApp
	messages   *message.MessageRepo
	chats      *chat.ChatApp
	tokenSvc   *auth.TokenService
}

func NewMessageApp(sessionApp *session.SessionApp,
	messages *message.MessageRepo,
	chats *chat.ChatApp,
	tokenSvc *auth.TokenService) *MessageApp {
	return &MessageApp{
		sessionApp: sessionApp,
		messages:   messages,
		chats:      chats,
		tokenSvc:   tokenSvc,
	}
}
