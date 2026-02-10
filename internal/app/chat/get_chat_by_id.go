package chat

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (c *ChatApp) GetChatById(tokenStr string, id int) (*domain.Chat, error) {
	session, err := c.sessionApp.GetSession(tokenStr)

	if err != nil {
		return nil, err
	}

	chat, err := c.chats.GetById(id)

	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, domain.NotFound("chat not found")
	}

	if !c.HasMember(chat, session.UserID) {
		return nil, domain.NotFound("chat not found")
	}

	return chat, nil
}
