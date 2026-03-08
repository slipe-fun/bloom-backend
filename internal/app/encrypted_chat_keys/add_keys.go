package encryptedchatkeys

import (
	"errors"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (k *EncryptedChatKeysApp) AddKeys(userID, recipientID, chatID int, keys []*domain.EncryptedChatKeys) ([]*domain.EncryptedChatKeys, int, error) {
	chat, err := k.chats.GetByID(chatID)
	if err != nil {
		logger.LogError(err.Error(), "encrypted-chat-keys-app")
		return nil, 0, domain.NotFound("failed to chat")
	}

	isUserInChat := false
	isRecipientInChat := false

	for _, member := range chat.Members {
		if member.ID == userID {
			isUserInChat = true
		}
		if member.ID == recipientID {
			isRecipientInChat = true
		}
	}

	if !isUserInChat {
		return nil, 0, errors.New("user is not a member of the chat")
	}

	if !isRecipientInChat {
		return nil, 0, errors.New("recipient is not a member of the chat")
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

	oldKeys, err := k.keys.GetBySessionIDsAndChatID(sessionIDs, chatID)
	if err != nil {
		return nil, 0, domain.Failed("failed to fetch existing keys")
	}

	if len(oldKeys) > 0 {
		oldIDs := make([]int, 0, len(oldKeys))
		for _, k := range oldKeys {
			oldIDs = append(oldIDs, k.ID)
		}

		err = k.keys.DeleteByIDs(oldIDs)
		if err != nil {
			return nil, 0, domain.Failed("failed to delete old keys")
		}
	}

	createdKeys, err := k.keys.Create(keys)
	if err != nil {
		return nil, 0, domain.Failed("failed to create encrypted chat keys")
	}

	return createdKeys, recipientID, nil
}
