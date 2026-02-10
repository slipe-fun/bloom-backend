package server

func (r *ServerMemberRepo) Delete(server_id, member_id int) error {
	query := `DELETE FROM server_members WHERE server_id = $1 AND member_id = $2`

	_, err := r.db.Exec(query, server_id, member_id)
	if err != nil {
		return err
	}

	return nil
}
