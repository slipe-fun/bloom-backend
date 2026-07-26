package message

import (
	"context"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (m *MessageApp) GetChatLastReadMessage(ctx context.Context, user_id, chatID int) (*domain.Message, error) {
	chat, err := m.chats.GetChatByID(ctx, user_id, chatID)
	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, domain.NotFound("chat not found")
	}

	if !m.chats.HasMember(ctx, chat, user_id) {
		return nil, domain.NotFound("chat not found")
	}

	message, err := m.messages.GetChatLastReadMessage(chat.ID)
	if err != nil {
		return nil, domain.NotFound("last read message not found")
	}

	return message, nil
}
