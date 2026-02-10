package friend

func (r *FriendRepo) Delete(userID int, friendID int) error {
	query := `
		DELETE FROM friends
		WHERE
		    (user_id = $1 AND friend_id = $2)
		 OR (user_id = $2 AND friend_id = $1)
	`

	_, err := r.db.Exec(query, userID, friendID)
	return err
}
