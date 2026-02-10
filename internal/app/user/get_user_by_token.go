package user

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (u *UserApp) GetUserByToken(tokenStr string) (*domain.User, error) {
	session, err := u.sessionApp.GetSession(tokenStr)
	if err != nil {
		return nil, err
	}

	user, err := u.users.GetById(session.UserID)
	if err != nil {
		logger.LogError(err.Error(), "user-app")
		return nil, domain.Failed("failed to get user")
	}

	return user, nil
}
