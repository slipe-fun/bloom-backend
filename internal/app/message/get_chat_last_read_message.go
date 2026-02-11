package message

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (m *MessageApp) GetChatLastReadMessage(token string, chatID int) (*domain.Message, error) {
	session, err := m.sessionApp.GetSession(token)
	if err != nil {
		return nil, err
	}

	chat, err := m.chats.GetChatByID(token, chatID)
	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, domain.NotFound("chat not found")
	}

	if !m.chats.HasMember(chat, session.UserID) {
		return nil, domain.NotFound("chat not found")
	}

	message, err := m.messages.GetChatLastReadMessage(chat.ID)
	if err != nil {
		return nil, domain.NotFound("last read message not found")
	}

	return message, nil
}
