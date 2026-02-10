package session

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (s *SessionApp) DeleteSession(id int, token string) error {
	userID, err := s.tokenSvc.ExtractUserID(token)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return domain.Unauthorized("failed to extract token")
	}

	_, err = s.users.GetById(userID)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return domain.NotFound("user not found")
	}

	session, err := s.session.GetById(id)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return domain.NotFound("session not found")
	}

	if session.UserID != userID {
		return domain.NotFound("session not found")
	}

	deleteSessionErr := s.session.Delete(id)
	if deleteSessionErr != nil {
		logger.LogError(deleteSessionErr.Error(), "session-app")
		return domain.Failed("failed to delete session")
	}

	return nil
}
