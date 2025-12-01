package UserApp

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (u *UserApp) GetUserByToken(tokenStr string) (*domain.User, error) {
	session, err := u.sessionApp.GetSession(tokenStr)
	if err != nil {
		return nil, err
	}

	user, err := u.users.GetById(session.UserID)
	if err != nil {
		return nil, domain.Failed("failed to get user")
	}

	return user, nil
}
