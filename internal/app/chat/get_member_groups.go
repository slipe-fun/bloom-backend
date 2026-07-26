package chat

import (
	"context"
	"database/sql"
	"errors"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (c *ChatApp) GetMemberGroups(ctx context.Context, memberID int) ([]domain.GroupMember, error) {
	groups, err := c.groups.GetByMemberID(ctx, memberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NotFound("groups not found")
		}
		return nil, domain.Failed("failed to get groups")
	}

	return groups, nil
}
