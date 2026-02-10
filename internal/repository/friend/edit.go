package friend

func (r *FriendRepo) EditStatus(userID, friendID int, status string) error {
	query := `
		UPDATE friends
		SET status = $1
		WHERE
		    (user_id = $2 AND friend_id = $3)
		 OR (user_id = $3 AND friend_id = $2)
	`

	_, err := r.db.Exec(query, status, userID, friendID)
	return err
}
