package verification

import (
	"errors"
)

func (r *VerificationRepo) DeleteByEmailAndCode(email string, code string) error {
	query := `DELETE FROM verification_codes WHERE email = $1 AND code = $2`

	result, err := r.db.Exec(query, email, code)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no deleted codes")
	}

	return nil
}
