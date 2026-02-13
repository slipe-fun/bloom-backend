package friend

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *FriendRepo) Create(friend *domain.FriendRow) (*domain.FriendRow, error) {
	query := `INSERT INTO friends (user_id, friend_id, status)
			  VALUES ($1, $2, $3)
			  RETURNING id, user_id, friend_id, status`

	var created domain.FriendRow

	start := time.Now()

	err := r.db.QueryRow(query, friend.UserID, friend.FriendID, friend.Status).
		Scan(&created.ID, &created.UserID, &created.FriendID, &created.Status)

	duration := time.Since(start)

	metrics.ObserveDB("friend_create", duration, err)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
