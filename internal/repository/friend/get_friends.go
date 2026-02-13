package friend

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *FriendRepo) GetFriends(userID int, status string, limit, offset int) ([]domain.Friend, error) {
	query := `SELECT
				id,
				CASE
					WHEN user_id = $1 THEN friend_id
					ELSE user_id
				END AS friend_id,
				status
			FROM friends
			WHERE
				status = $2
				AND (user_id = $1 OR friend_id = $1)
			ORDER BY id
			LIMIT $3 OFFSET $4;
			`

	var friends []domain.Friend

	start := time.Now()

	err := r.db.Select(&friends, query, userID, status, limit, offset)

	duration := time.Since(start)

	metrics.ObserveDB("friends_get", duration, err)

	if err != nil {
		return nil, err
	}

	return friends, nil
}
