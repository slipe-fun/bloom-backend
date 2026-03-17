package encryptedchatkeys

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

type chatSessionPair struct {
	ChatID    int
	SessionID int
}

func (k *EncryptedChatKeysApp) AddKeys(userID int, batch []domain.AddKeyBatchItem) ([]*domain.EncryptedChatKeys, error) {
	if len(batch) == 0 {
		return nil, nil
	}

	uniqueChatIDsMap := make(map[int]struct{})
	uniqueSessionIDsMap := make(map[int]struct{})
	seenPairs := make(map[chatSessionPair]struct{})

	var chatIDs []int
	var sessionIDs []int
	var deduplicatedBatch []domain.AddKeyBatchItem

	for _, item := range batch {
		pair := chatSessionPair{ChatID: item.ChatID, SessionID: item.Key.SessionID}
		if _, exists := seenPairs[pair]; exists {
			continue
		}
		seenPairs[pair] = struct{}{}
		deduplicatedBatch = append(deduplicatedBatch, item)

		if _, exists := uniqueChatIDsMap[item.ChatID]; !exists {
			uniqueChatIDsMap[item.ChatID] = struct{}{}
			chatIDs = append(chatIDs, item.ChatID)
		}
		if _, exists := uniqueSessionIDsMap[item.Key.SessionID]; !exists {
			uniqueSessionIDsMap[item.Key.SessionID] = struct{}{}
			sessionIDs = append(sessionIDs, item.Key.SessionID)
		}
	}

	chats, err := k.chats.GetByIDs(chatIDs)
	if err != nil {
		logger.LogError(err.Error(), "encrypted-chat-keys-app")
		return nil, domain.Failed("failed to fetch chats")
	}

	chatMembersCache := make(map[int]map[int]bool)
	for _, chat := range chats {
		membersMap := make(map[int]bool, len(chat.Members))
		for _, member := range chat.Members {
			membersMap[member.ID] = true
		}
		chatMembersCache[chat.ID] = membersMap
	}

	sessions, err := k.session.GetSessionByIDs(sessionIDs)
	if err != nil {
		return nil, domain.Failed("failed to fetch sessions")
	}

	sessionOwnerMap := make(map[int]int)
	for _, s := range sessions {
		sessionOwnerMap[s.ID] = s.UserID
	}

	validKeysToInsert := make([]*domain.EncryptedChatKeys, 0, len(deduplicatedBatch))
	validPairs := make(map[chatSessionPair]struct{})

	for _, item := range deduplicatedBatch {
		members, chatExists := chatMembersCache[item.ChatID]
		if !chatExists {
			continue
		}
		if !members[userID] {
			continue
		}
		if !members[item.RecipientID] {
			continue
		}

		sessionOwnerID, sessionExists := sessionOwnerMap[item.Key.SessionID]
		if !sessionExists {
			continue
		}
		if sessionOwnerID != item.RecipientID {
			continue
		}

		item.Key.ChatID = item.ChatID
		validKeysToInsert = append(validKeysToInsert, item.Key)
		validPairs[chatSessionPair{ChatID: item.ChatID, SessionID: item.Key.SessionID}] = struct{}{}
	}

	if len(validKeysToInsert) == 0 {
		return []*domain.EncryptedChatKeys{}, nil
	}

	oldKeys, err := k.keys.GetBySessionIDsAndChatIDs(sessionIDs, chatIDs)
	if err != nil {
		return nil, domain.Failed("failed to fetch existing keys")
	}

	var oldIDsToDelete []int
	for _, oldKey := range oldKeys {
		pair := chatSessionPair{ChatID: oldKey.ChatID, SessionID: oldKey.SessionID}
		if _, isBeingReplaced := validPairs[pair]; isBeingReplaced {
			oldIDsToDelete = append(oldIDsToDelete, oldKey.ID)
		}
	}

	if len(oldIDsToDelete) > 0 {
		err = k.keys.DeleteByIDs(oldIDsToDelete)
		if err != nil {
			return nil, domain.Failed("failed to delete old keys")
		}
	}

	createdKeys, err := k.keys.Create(validKeysToInsert)
	if err != nil {
		return nil, domain.Failed("failed to create encrypted chat keys")
	}

	return createdKeys, nil
}
