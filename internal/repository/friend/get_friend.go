package friend

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *FriendRepo) GetFriend(userID int, friendID int) (*domain.FriendRow, error) {
	query := `
		SELECT
			id,
			user_id,
			friend_id,
			status
		FROM friends
		WHERE
			(user_id = $1 AND friend_id = $2)
			OR (user_id = $2 AND friend_id = $1)
		LIMIT 1
	`

	var friend domain.FriendRow

	start := time.Now()

	err := r.db.Get(&friend, query, userID, friendID)

	duration := time.Since(start)

	metrics.ObserveDB("friend_get", duration, err)

	if err != nil {
		return nil, err
	}

	return &friend, nil
}
