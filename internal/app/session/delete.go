package SessionApp

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (s *SessionApp) DeleteSession(id int, token string) error {
	userID, err := s.tokenSvc.ExtractUserID(token)
	if err != nil {
		return domain.Unauthorized("failed to extract token")
	}

	_, err = s.users.GetById(userID)
	if err != nil {
		return domain.NotFound("user not found")
	}

	session, err := s.session.GetById(id)
	if err != nil {
		return domain.NotFound("session not found")
	}

	if session.UserID != userID {
		return domain.NotFound("session not found")
	}

	deleteSessionErr := s.session.Delete(id)
	if deleteSessionErr != nil {
		return domain.Failed("failed to delete session")
	}

	return nil
}
