package server

func (r *ServerMessageRepo) Edit(id int, newContent string) error {
	query := `UPDATE server_messages SET content = $1 WHERE id = $2`

	_, err := r.db.Exec(query, newContent, id)
	if err != nil {
		return err
	}

	return nil
}
