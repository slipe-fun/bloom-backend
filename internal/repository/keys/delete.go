package keys

func (k *KeysRepo) Delete(user_id int) error {
	query := `DELETE FROM keys WHERE user_id = $1`

	_, err := k.db.Exec(query, user_id)
	if err != nil {
		return err
	}

	return nil
}
