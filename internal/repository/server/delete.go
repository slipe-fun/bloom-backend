package server

func (r *ServerRepo) Delete(id int) error {
	query := `DELETE FROM servers WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
