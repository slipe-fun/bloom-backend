package session

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (s *SessionApp) GetSessionByUserIDs(ids []int) ([]*domain.Session, error) {
	sessions, err := s.session.GetByUserIDs(ids)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return nil, domain.Failed("failed to get sessions")
	}

	return sessions, nil
}
