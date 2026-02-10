package session

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/crypto"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (s *SessionApp) CreateSession(user_id int) (string, error) {
	user, err := s.users.GetById(user_id)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return "", domain.NotFound("user not found")
	}

	token, err := s.jwtSvc.GenerateToken(user.ID)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return "", domain.Failed("failed to generate token")
	}

	_, err = s.session.Create(&domain.Session{
		Token:  crypto.HashSHA256(token),
		UserID: user.ID,
	})
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return "", domain.Failed("failed to generate token")
	}

	return token, nil
}
