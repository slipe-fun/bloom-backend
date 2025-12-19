package ChatApp

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/service/logger"
)

func (c *ChatApp) GetChatsByUserId(tokenStr string) ([]*domain.ChatWithLastMessage, error) {
	session, err := c.sessionApp.GetSession(tokenStr)
	if err != nil {
		return nil, err
	}

	chats, err := c.chats.GetByUserId(session.UserID)

	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, domain.NotFound("chats not found")
	}

	return chats, nil
}
