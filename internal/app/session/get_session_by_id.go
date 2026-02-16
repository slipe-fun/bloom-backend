package session

import "github.com/slipe-fun/skid-backend/internal/domain"

func (s *SessionApp) GetSessionByID(id int) (*domain.Session, error) {
	session, err := s.session.GetByID(id)
	if err != nil {
		return nil, domain.Failed("failed to get session")
	}

	return session, err
}
