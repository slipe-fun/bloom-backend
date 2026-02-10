package friend

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (f *FriendApp) GetFriendCount(token string) (int, error) {
	session, err := f.sessionApp.GetSession(token)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return 0, domain.Unauthorized("session not found")
	}

	count, err := f.friends.GetFriendCount(session.UserID)
	if err != nil {
		logger.LogError(err.Error(), "friend-app")
		return 0, domain.Failed("failed to get friend count")
	}

	return count, nil
}
