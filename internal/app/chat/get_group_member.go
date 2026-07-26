package chat

import (
	"context"
	"database/sql"
	"errors"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (c *ChatApp) GetGroupMember(ctx context.Context, groupID, memberID int) (*domain.GroupMember, error) {
	member, err := c.groups.GetByMemberAndChatID(ctx, groupID, memberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NotFound("member not found")
		}
		return nil, domain.Failed("failed to get group members")
	}

	return member, nil
}
