package MessageApp

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (m *MessageApp) UpdateMessagesSeenStatus(messageIDs []int, seenAt time.Time) error {
	updateMessageError := m.messages.UpdateMessagesSeenStatus(messageIDs, seenAt)
	if updateMessageError != nil {
		return domain.Failed("failed to update message seen status")
	}

	return nil
}
