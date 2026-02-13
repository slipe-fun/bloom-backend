package friend

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *FriendRepo) Delete(userID int, friendID int) error {
	query := `
		DELETE FROM friends
		WHERE
		    (user_id = $1 AND friend_id = $2)
		 OR (user_id = $2 AND friend_id = $1)
	`

	start := time.Now()

	_, err := r.db.Exec(query, userID, friendID)

	duration := time.Since(start)

	metrics.ObserveDB("friend_delete", duration, err)

	return err
}
