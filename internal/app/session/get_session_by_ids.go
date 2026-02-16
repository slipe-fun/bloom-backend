package session

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (s *SessionApp) GetSessionByIDs(ids []int) ([]*domain.Session, error) {
	sessions, err := s.session.GetByIDs(ids)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return nil, domain.Failed("failed to get sessions")
	}

	return sessions, nil
}
