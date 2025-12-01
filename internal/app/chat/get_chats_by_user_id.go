package ChatApp

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (c *ChatApp) GetChatsByUserId(tokenStr string) ([]*domain.Chat, error) {
	session, err := c.sessionApp.GetSession(tokenStr)
	if err != nil {
		return nil, err
	}

	chats, err := c.chats.GetByUserId(session.UserID)

	if err != nil {
		return nil, domain.NotFound("chats not found")
	}

	return chats, nil
}
