package chat

import (
	"context"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (c *ChatApp) CreatePrivateChat(user_id, recipient int, handshake domain.Handshake) (*domain.Chat, error) {
	chat, err := c.chats.Create(&domain.RawChat{
		Members: []domain.Member{
			{
				ID: user_id,
			},
			{
				ID: recipient,
			},
		},
		Type:      "private",
		Handshake: &handshake,
	})

	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, domain.Failed("failed to create chat")
	}

	return chat, nil
}

func (c *ChatApp) CreateGroupChat(ctx context.Context, title string, members []domain.GroupMember, invitedByID int) (*domain.Chat, error) {
	chat, err := c.chats.Create(&domain.RawChat{
		Type:  "group",
		Title: title,
	})

	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, domain.Failed("failed to create chat")
	}

	err = c.groups.CreateMany(ctx, chat.ID, invitedByID, members)
	if err != nil {
		return nil, err
	}

	return chat, nil
}
