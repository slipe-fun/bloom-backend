package MessageApp

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (c *MessageApp) CreateMessage(tokenStr string, chatId int, message *domain.Message) (*domain.Message, error) {
	_, err := c.sessionApp.GetSession(tokenStr)
	if err != nil {
		return nil, err
	}

	_, chatErr := c.chats.GetChatById(tokenStr, chatId)
	if chatErr != nil {
		return nil, chatErr
	}

	message, messageErr := c.messages.Create(message)
	if messageErr != nil {
		return nil, domain.Failed("failed to create message")
	}

	return message, nil
}
