package chat

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (c *ChatApp) AddKeys(user_id int, chat *domain.Chat, kyberPublicKey string, ecdhPublicKey string, edPublicKey string) error {
	for i := range chat.Members {
		if chat.Members[i].ID == user_id {
			chat.Members[i].KyberPublicKey = kyberPublicKey
			chat.Members[i].EcdhPublicKey = ecdhPublicKey
			chat.Members[i].EdPublicKey = edPublicKey
		}
	}

	updateChatError := c.chats.UpdateChat(chat)
	if updateChatError != nil {
		logger.LogError(updateChatError.Error(), "chat-app")
		return domain.Failed("failed to update chat keys")
	}

	return nil
}
