package friend

import (
	"database/sql"
	"errors"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (f *FriendApp) SendRequest(token string, receiverID int) error {
	session, err := f.sessionApp.GetSession(token)
	if err != nil {
		return err
	}

	if session.UserID == receiverID {
		return domain.InvalidData("cannot add yourself as friend")
	}

	receiver, err := f.users.GetById(receiverID)
	if err != nil {
		logger.LogError(err.Error(), "user-app")
		return domain.NotFound("receiver not found")
	}

	friend, err := f.friends.GetFriend(session.UserID, receiver.ID)
	if err == nil {
		if friend.Status == "pending" {
			if friend.FriendID == session.UserID {
				if err := f.friends.EditStatus(session.UserID, receiverID, "accepted"); err != nil {
					logger.LogError(err.Error(), "friend-app")
					return domain.Failed("failed to accept friend request")
				}
				return nil
			}

			return domain.AlreadyExists("friend request already sent")
		}

		if friend.Status == "accepted" {
			return domain.AlreadyExists("already friends")
		}

		return domain.AlreadyExists("invalid friend state")
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.LogError(err.Error(), "friend-app")
		return domain.Failed("failed to check friend relation")
	}

	_, err = f.friends.Create(&domain.FriendRow{
		UserID:   session.UserID,
		FriendID: receiver.ID,
		Status:   "pending",
	})
	if err != nil {
		logger.LogError(err.Error(), "user-app")
		return domain.Failed("failed to send friend request")
	}

	return nil
}
