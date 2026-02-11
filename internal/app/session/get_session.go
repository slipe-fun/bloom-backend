package session

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/crypto"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (s *SessionApp) GetSession(token string) (*domain.Session, error) {
	userID, err := s.tokenSvc.ExtractUserID(token)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return nil, domain.Unauthorized("failed to extract token")
	}

	_, err = s.users.GetByID(userID)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return nil, domain.Unauthorized("user not found")
	}

	session, err := s.session.GetByToken(crypto.HashSHA256(token))
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return nil, domain.Unauthorized("session not found")
	}

	if session.UserID != userID {
		return nil, domain.Unauthorized("session not found")
	}

	return session, nil
}
