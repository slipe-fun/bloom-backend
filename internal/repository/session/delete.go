package session

func (r *SessionRepo) Delete(id int) error {
	query := `DELETE FROM sessions WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
