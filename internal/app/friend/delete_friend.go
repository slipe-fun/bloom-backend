package friend

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (f *FriendApp) DeleteFriend(token string, friend_id int) error {
	session, err := f.sessionApp.GetSession(token)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return domain.NotFound("session not found")
	}

	_, err = f.friends.GetFriend(session.UserID, friend_id)
	if err != nil {
		logger.LogError(err.Error(), "friend-app")
		return domain.NotFound("friend not found")
	}

	err = f.friends.Delete(session.UserID, friend_id)
	if err != nil {
		logger.LogError(err.Error(), "friend-app")
		return domain.Failed("failed to delete friend")
	}

	return nil
}
