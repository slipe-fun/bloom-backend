package session

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *SessionRepo) GetByUserID(id int) ([]*domain.Session, error) {
	var sessions []*domain.Session

	query := `SELECT id, token, user_id, created_at FROM sessions WHERE user_id = $1`
	err := r.db.Select(&sessions, query, id)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}
