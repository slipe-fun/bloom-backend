package session

import (
	"time"

	"github.com/lib/pq"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *SessionRepo) GetByIDs(ids []int) ([]*domain.Session, error) {
	query := `
		SELECT id, token, user_id, identity_pub, ecdh_pub, kyber_pub, revoked_at, created_at
		FROM sessions
		WHERE id = ANY($1)
	`

	start := time.Now()

	var sessions []*domain.Session
	err := r.db.Select(&sessions, query, pq.Array(ids))

	duration := time.Since(start)
	metrics.ObserveDB("session_get_many", duration, err)

	if err != nil {
		return nil, err
	}

	return sessions, nil
}
