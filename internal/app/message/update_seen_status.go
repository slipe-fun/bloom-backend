package message

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (m *MessageApp) UpdateMessagesSeenStatus(user_id, chatID int, messageIDs []int) (*[]int, *time.Time, *domain.Chat, error) {
	chat, err := m.chats.GetChatByID(user_id, chatID)
	if err != nil {
		return nil, nil, nil, err
	}

	if len(messageIDs) == 0 {
		return nil, nil, nil, domain.InvalidData("no message ids provided")
	}

	var validMessages []int
	for i := range messageIDs {
		message, err := m.messages.GetByID(messageIDs[i])
		if err != nil {
			continue
		}

		if message.Seen != nil || message.ChatID == 0 || message.ChatID != chat.ID {
			continue
		}

		validMessages = append(validMessages, messageIDs[i])
	}

	if len(validMessages) == 0 {
		return nil, nil, nil, domain.InvalidData("no valid messages to update")
	}

	seenAt := time.Now()
	updateMessageError := m.messages.UpdateMessagesSeenStatus(validMessages, seenAt)
	if updateMessageError != nil {
		logger.LogError(updateMessageError.Error(), "message-app")
		return nil, nil, nil, domain.Failed("failed to update message seen status")
	}

	return &validMessages, &seenAt, chat, nil
}
