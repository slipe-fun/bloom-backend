package session

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *SessionRepo) GetByID(id int) (*domain.Session, error) {
	var session domain.Session

	query := `SELECT id, token, user_id, created_at FROM sessions WHERE id = $1`

	start := time.Now()

	err := r.db.Get(&session, query, id)

	duration := time.Since(start)

	metrics.ObserveDB("session_get_by_id", duration, err)

	if err != nil {
		return nil, err
	}

	return &session, nil
}
