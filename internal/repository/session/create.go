package session

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *SessionRepo) Create(session *domain.Session) (*domain.Session, error) {
	query := `INSERT INTO sessions (token, user_id) 
	          VALUES ($1, $2) 
	          RETURNING id, token, user_id, created_at`

	var created domain.Session

	start := time.Now()

	err := r.db.QueryRow(query, session.Token, session.UserID).
		Scan(&created.ID, &created.Token, &created.UserID, &created.CreatedAt)

	duration := time.Since(start)

	metrics.ObserveDB("session_create", duration, err)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
