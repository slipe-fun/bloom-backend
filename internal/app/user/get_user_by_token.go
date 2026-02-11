package user

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (u *UserApp) GetUserByToken(token string) (*domain.User, error) {
	session, err := u.sessionApp.GetSession(token)
	if err != nil {
		return nil, err
	}

	user, err := u.users.GetByID(session.UserID)
	if err != nil {
		logger.LogError(err.Error(), "user-app")
		return nil, domain.Failed("failed to get user")
	}

	return user, nil
}
