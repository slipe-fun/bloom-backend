package SessionApp

import (
	"errors"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/service"
)

func (s *SessionApp) GetUserSessions(token string) ([]*domain.Session, error) {
	session, err := s.session.GetByToken(service.HashSHA256(token))
	if err != nil {
		return nil, errors.New("session not found")
	}

	user, err := s.users.GetById(session.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	sessions, err := s.session.GetByUserId(user.ID)
	if err != nil {
		return nil, errors.New("failed to get user sessions")
	}

	return sessions, nil
}
