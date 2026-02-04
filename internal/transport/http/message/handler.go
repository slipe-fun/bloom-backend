package message

import (
	ChatApp "github.com/slipe-fun/skid-backend/internal/app/chat"
	MessageApp "github.com/slipe-fun/skid-backend/internal/app/message"
	UserApp "github.com/slipe-fun/skid-backend/internal/app/user"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

type MessageHandler struct {
	chatApp    *ChatApp.ChatApp
	userApp    *UserApp.UserApp
	messageApp *MessageApp.MessageApp
	wsHub      *types.Hub
}

func NewMessageHandler(chatApp *ChatApp.ChatApp, userApp *UserApp.UserApp, messageApp *MessageApp.MessageApp, wsHub *types.Hub) *MessageHandler {
	return &MessageHandler{
		chatApp:    chatApp,
		userApp:    userApp,
		messageApp: messageApp,
		wsHub:      wsHub,
	}
}
