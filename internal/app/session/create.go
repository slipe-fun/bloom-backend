package SessionApp

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/service"
)

func (s *SessionApp) CreateSession(user_id int) (string, error) {
	user, err := s.users.GetById(user_id)
	if err != nil {
		return "", domain.NotFound("user not found")
	}

	token, err := s.jwtSvc.GenerateToken(user.ID)
	if err != nil {
		return "", domain.Failed("failed to generate token")
	}

	_, err = s.session.Create(&domain.Session{
		Token:  service.HashSHA256(token),
		UserID: user.ID,
	})
	if err != nil {
		return "", domain.Failed("failed to generate token")
	}

	return token, nil
}
