package FriendRepo

func (r *FriendRepo) Delete(userID int, friendID int) error {
	query := `DELETE FROM friends WHERE user_id = $1 AND friend_id = $2`

	_, err := r.db.Exec(query, userID, friendID)
	if err != nil {
		return err
	}

	return nil
}
