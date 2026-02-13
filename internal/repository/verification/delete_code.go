package verification

import (
	"errors"
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *VerificationRepo) DeleteByEmailAndCode(email string, code string) error {
	query := `DELETE FROM verification_codes WHERE email = $1 AND code = $2`

	start := time.Now()

	result, err := r.db.Exec(query, email, code)

	duration := time.Since(start)

	metrics.ObserveDB("verification_code_delete_by_email_and_code", duration, err)

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
