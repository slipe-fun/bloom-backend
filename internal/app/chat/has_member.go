package chat

import (
	"context"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (c *ChatApp) HasMember(ctx context.Context, chat *domain.Chat, memberID int) bool {
	switch chat.Type {
	case "private":
		for _, m := range chat.Members {
			if m.ID == memberID {
				return true
			}
		}
	case "group":
		_, err := c.groups.GetByMemberAndChatID(ctx, memberID, chat.ID)
		if err != nil {
			return false
		}
		return true
	default:
		return false
	}
	return false
}
