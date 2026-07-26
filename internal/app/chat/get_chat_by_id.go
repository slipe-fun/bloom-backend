package chat

import (
	"context"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (c *ChatApp) GetChatByID(ctx context.Context, user_id int, id int) (*domain.Chat, error) {
	chat, err := c.chats.GetByID(id)
	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, domain.NotFound("chat not found")
	}

	if !c.HasMember(ctx, chat, user_id) {
		return nil, domain.NotFound("chat not found")
	}

	return chat, nil
}
