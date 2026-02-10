package friend

import "github.com/slipe-fun/skid-backend/internal/domain"

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
	err := r.db.Get(&friend, query, userID, friendID)
	if err != nil {
		return nil, err
	}

	return &friend, nil
}
