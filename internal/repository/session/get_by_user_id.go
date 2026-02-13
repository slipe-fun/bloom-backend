package session

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *SessionRepo) GetByUserID(id int) ([]*domain.Session, error) {
	var sessions []*domain.Session

	query := `SELECT id, token, user_id, created_at FROM sessions WHERE user_id = $1`

	start := time.Now()

	err := r.db.Select(&sessions, query, id)

	duration := time.Since(start)

	metrics.ObserveDB("sessions_get_by_user_id", duration, err)

	if err != nil {
		return nil, err
	}

	return sessions, nil
}
