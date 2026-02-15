package session

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/crypto"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (s *SessionApp) CreateSession(user_id int) (string, *domain.Session, error) {
	user, err := s.users.GetByID(user_id)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return "", nil, domain.NotFound("user not found")
	}

	token, err := s.jwtSvc.GenerateToken(user.ID)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return "", nil, domain.Failed("failed to generate token")
	}

	session, err := s.session.Create(&domain.Session{
		Token:  crypto.HashSHA256(token),
		UserID: user.ID,
	})
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return "", nil, domain.Failed("failed to generate token")
	}

	return token, session, nil
}
