package chat

import (
	"context"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (c *ChatApp) InviteUsersToGroup(ctx context.Context, chatID, invitedByID int, members []domain.GroupMember) error {
	err := c.groups.CreateMany(ctx, chatID, invitedByID, members)
	if err != nil {
		return domain.Failed("failed to add members")
	}

	return nil
}
