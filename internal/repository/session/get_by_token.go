package session

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *SessionRepo) GetByToken(token string) (*domain.Session, error) {
	var session domain.Session

	query := `SELECT id, token, user_id, created_at FROM sessions WHERE token = $1`
	err := r.db.Get(&session, query, token)

	if err != nil {
		return nil, err
	}

	return &session, nil
}
