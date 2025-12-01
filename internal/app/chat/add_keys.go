package ChatApp

import "github.com/slipe-fun/skid-backend/internal/domain"

func (c *ChatApp) AddKeys(tokenStr string, chat *domain.Chat, kyberPublicKey string, ecdhPublicKey string, edPublicKey string) error {
	session, err := c.sessionApp.GetSession(tokenStr)
	if err != nil {
		return err
	}

	for i := range chat.Members {
		if chat.Members[i].ID == session.UserID {
			chat.Members[i].KyberPublicKey = kyberPublicKey
			chat.Members[i].EcdhPublicKey = ecdhPublicKey
			chat.Members[i].EdPublicKey = edPublicKey
		}
	}

	updateChatError := c.chats.UpdateChat(chat)
	if updateChatError != nil {
		return domain.Failed("failed to update chat keys")
	}

	return nil
}
