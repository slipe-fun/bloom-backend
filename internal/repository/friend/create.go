package FriendRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *FriendRepo) Create(friend *domain.Friend) (*domain.Friend, error) {
	query := `INSERT INTO friends (user_id, friend_id, status)
			  VALUES ($1, $2, $3)
			  RETURNING id, user_id, friend_id, status`

	var created domain.Friend
	err := r.db.QueryRow(query, friend.UserID, friend.FriendID, friend.Status).
		Scan(&created.ID, &created.UserID, &created.FriendID, &created.Status)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
