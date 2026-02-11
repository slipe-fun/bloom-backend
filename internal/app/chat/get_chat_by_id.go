package chat

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (c *ChatApp) GetChatByID(token string, id int) (*domain.Chat, error) {
	session, err := c.sessionApp.GetSession(token)

	if err != nil {
		return nil, err
	}

	chat, err := c.chats.GetByID(id)

	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, domain.NotFound("chat not found")
	}

	if !c.HasMember(chat, session.UserID) {
		return nil, domain.NotFound("chat not found")
	}

	return chat, nil
}
