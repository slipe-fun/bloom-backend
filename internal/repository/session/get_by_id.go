package session

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *SessionRepo) GetById(id int) (*domain.Session, error) {
	var session domain.Session

	query := `SELECT id, token, user_id, created_at FROM sessions WHERE id = $1`
	err := r.db.Get(&session, query, id)

	if err != nil {
		return nil, err
	}

	return &session, nil
}
