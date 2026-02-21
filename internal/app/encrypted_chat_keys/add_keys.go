package encryptedchatkeys

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (k *EncryptedChatKeysApp) AddKeys(userID, chatID int, keys []*domain.EncryptedChatKeys) ([]*domain.EncryptedChatKeys, int, error) {
	chat, err := k.chats.GetByID(chatID)
	if err != nil {
		logger.LogError(err.Error(), "encrypted-chat-keys-app")
		return nil, 0, domain.NotFound("failed to chat")
	}

	recipientID := 0
	isUserInChat := false

	for _, member := range chat.Members {
		if member.ID == userID {
			isUserInChat = true
		} else {
			recipientID = member.ID
		}
	}

	if !isUserInChat {
		return nil, 0, domain.NotFound("failed to chat")
	}

	if recipientID == 0 {
		return nil, 0, domain.InvalidData("recipient is not member of the provided chat")
	}

	seen := make(map[int]struct{}, len(keys))
	sessionIDs := make([]int, 0, len(keys))

	for _, key := range keys {
		if _, exists := seen[key.SessionID]; exists {
			continue
		}
		seen[key.SessionID] = struct{}{}
		sessionIDs = append(sessionIDs, key.SessionID)
	}

	sessions, err := k.session.GetSessionByIDs(sessionIDs)
	if err != nil {
		return nil, 0, domain.Failed("failed to get provided sessions")
	}

	if len(sessions) != len(sessionIDs) {
		return nil, 0, domain.InvalidData("one or more sessions not found")
	}

	for _, s := range sessions {
		if s.UserID != recipientID {
			return nil, 0, domain.InvalidData("one of provided sessions ids is not belongs to recipient")
		}
	}

	createdKeys, err := k.keys.Create(keys)
	if err != nil {
		return nil, 0, domain.Failed("failed to create encrypted chat keys")
	}

	return createdKeys, recipientID, nil
}
