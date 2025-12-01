package ChatApp

import "github.com/slipe-fun/skid-backend/internal/domain"

func (c *ChatApp) GetChatWithUsers(tokenStr string, recipient int) (*domain.Chat, error) {
	session, err := c.sessionApp.GetSession(tokenStr)
	if err != nil {
		return nil, err
	}

	chat, err := c.chats.GetWithUsers(session.UserID, recipient)

	if err != nil {
		return nil, domain.NotExpired("chats not found")
	}

	return chat, nil
}
