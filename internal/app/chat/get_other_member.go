package chat

import "github.com/slipe-fun/skid-backend/internal/domain"

func (c *ChatApp) GetOtherMember(chat *domain.Chat, memberID int) *domain.User {
	for _, m := range chat.Members {
		if m.ID != memberID {
			return &m
		}
	}
	return nil
}
