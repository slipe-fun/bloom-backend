package friend

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (f *FriendApp) GetFriends(token, status string, limit, offset int) ([]domain.Friend, error) {
	session, err := f.sessionApp.GetSession(token)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return nil, domain.Unauthorized("session not found")
	}

	friends, err := f.friends.GetFriends(session.UserID, status, limit, offset)
	if err != nil {
		logger.LogError(err.Error(), "friend-app")
		return nil, domain.Failed("failed to get friends")
	}

	return friends, nil
}
