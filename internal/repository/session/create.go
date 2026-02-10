package session

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *SessionRepo) Create(session *domain.Session) (*domain.Session, error) {
	query := `INSERT INTO sessions (token, user_id) 
	          VALUES ($1, $2) 
	          RETURNING id, token, user_id, created_at`

	var created domain.Session
	err := r.db.QueryRow(query, session.Token, session.UserID).
		Scan(&created.ID, &created.Token, &created.UserID, &created.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
