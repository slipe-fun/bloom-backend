package session

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *SessionRepo) AddKeys(id int, identity_pub, ecdh_pub, kyber_pub string) error {
	query := `
		UPDATE sessions
		SET identity_pub = $2,
			ecdh_pub = $3,
			kyber_pub = $4
		WHERE id = $1
	`

	start := time.Now()

	_, err := r.db.Exec(query, id, identity_pub, ecdh_pub, kyber_pub)

	duration := time.Since(start)

	metrics.ObserveDB("session_add_keys", duration, err)

	return err
}
