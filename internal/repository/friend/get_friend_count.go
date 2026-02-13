package friend

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *FriendRepo) GetFriendCount(userID int) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM friends
		WHERE
			(user_id = $1 OR friend_id = $1)
			AND status = 'accepted'
	`

	var count int

	start := time.Now()

	err := r.db.Get(&count, query, userID)

	duration := time.Since(start)

	metrics.ObserveDB("friend_get_count", duration, err)

	if err != nil {
		return 0, err
	}

	return count, nil
}
