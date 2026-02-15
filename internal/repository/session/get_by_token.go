package session

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *SessionRepo) GetByToken(token string) (*domain.Session, error) {
	var session domain.Session

	query := `SELECT id, token, user_id, identity_pub, ecdh_pub, kyber_pub, revoked_at, created_at FROM sessions WHERE token = $1`

	start := time.Now()

	err := r.db.Get(&session, query, token)

	duration := time.Since(start)

	metrics.ObserveDB("session_get_by_token", duration, err)

	if err != nil {
		return nil, err
	}

	return &session, nil
}
