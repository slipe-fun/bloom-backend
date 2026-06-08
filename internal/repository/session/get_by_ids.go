package session

import (
	"time"

	"github.com/lib/pq"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *SessionRepo) GetByIDs(ids []int) ([]*domain.Session, error) {
	var sessions []*domain.Session

	query := `
	SELECT s.id, s.token, s.user_id, s.revoked_at, s.created_at, u.public_id AS user_public_id
	FROM sessions s
	JOIN users u ON u.id = s.user_id
	WHERE s.id = ANY($1)
	`

	start := time.Now()

	err := r.db.Select(&sessions, query, pq.Array(ids))

	duration := time.Since(start)

	metrics.ObserveDB("session_get_many", duration, err)

	if err != nil {
		return nil, err
	}

	return sessions, nil
}
