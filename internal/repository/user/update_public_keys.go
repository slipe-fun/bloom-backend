package user

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *UserRepo) UpdatePublicKeys(userID int, kyber, ecdh, ed string) error {
	query := `
        UPDATE users
        SET kyber_public_key = $1,
            ecdh_public_key = $2,
            ed_public_key = $3
        WHERE id = $4
    `

	start := time.Now()

	_, err := r.db.Exec(query, kyber, ecdh, ed, userID)

	duration := time.Since(start)

	metrics.ObserveDB("user_update_public_keys", duration, err)

	if err != nil {
		return err
	}

	return nil
}
