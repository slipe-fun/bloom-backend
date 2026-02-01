package ServerChannelRepo

func (r *ServerChannelRepo) Delete(id int) error {
	query := `DELETE FROM server_channels WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
