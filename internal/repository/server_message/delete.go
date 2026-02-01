package ServerMessageRepo

func (r *ServerMessageRepo) Delete(id int) error {
	query := `DELETE FROM server_messages WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
