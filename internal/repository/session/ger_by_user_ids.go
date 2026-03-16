package session

import (
	"time"

	"github.com/lib/pq"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *SessionRepo) GetByUserIDs(ids []int) ([]*domain.Session, error) {
	var sessions []*domain.Session

	query := `
	SELECT id, token, user_id, identity_pub, ecdh_pub, kyber_pub, revoked_at, created_at
	FROM sessions
	WHERE user_id = ANY($1)
	`

	start := time.Now()

	err := r.db.Select(&sessions, query, pq.Array(ids))

	duration := time.Since(start)

	metrics.ObserveDB("sessions_get_by_user_ids", duration, err)

	if err != nil {
		return nil, err
	}

	return sessions, nil
}
