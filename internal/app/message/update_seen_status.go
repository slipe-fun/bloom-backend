package MessageApp

import "time"

func (m *MessageApp) UpdateMessagesSeenStatus(messageIDs []int, seenAt time.Time) error {
	return m.messages.UpdateMessagesSeenStatus(messageIDs, seenAt)
}
