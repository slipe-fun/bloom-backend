package FriendRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *FriendRepo) EditStatus(friend *domain.Friend) error {
	query := `UPDATE friends
		SET status = $1
		WHERE id = $2
		RETURNING id, user_id, friend_id, status
	`

	_, err := r.db.Exec(query, friend.Status, friend.ID)

	if err != nil {
		return err
	}

	return nil
}
