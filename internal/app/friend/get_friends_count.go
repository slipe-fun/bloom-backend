package friend

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (f *FriendApp) GetFriendCount(user_id int) (int, error) {
	count, err := f.friends.GetFriendCount(user_id)
	if err != nil {
		logger.LogError(err.Error(), "friend-app")
		return 0, domain.Failed("failed to get friend count")
	}

	return count, nil
}
