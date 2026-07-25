package user

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (u *UserApp) GetUsersByPublicIDs(ids []string) ([]domain.User, error) {
	if len(ids) == 0 {
		return nil, domain.NotFound("users not found")
	}

	users, err := u.users.GetUsersByPublicIDs(ids)
	if err != nil {
		logger.LogError(err.Error(), "user-app")
		return nil, domain.NotFound("users not found")
	}

	if len(users) == 0 {
		return nil, domain.NotFound("users not found")
	}

	return users, nil
}
