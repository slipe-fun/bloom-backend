package user

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (u *UserApp) SearchUsersByUsername(username string, limit, offset int) ([]*domain.User, error) {
	users, err := u.users.SearchUsersByUsername(username, limit, offset)
	if err != nil {
		logger.LogError(err.Error(), "user-app")
		return nil, domain.NotFound("users not found")
	}

	if len(users) == 0 {
		return nil, domain.NotFound("users not found")
	}

	return users, nil
}
