package friend

func (r *FriendRepo) GetFriendCount(userID int) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM friends
		WHERE
			(user_id = $1 OR friend_id = $1)
			AND status = 'accepted'
	`

	var count int
	err := r.db.Get(&count, query, userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
