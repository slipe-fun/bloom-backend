package user

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (u *UserApp) GetUsersByIDs(ids []int) ([]domain.User, error) {
	user, err := u.users.GetUsersByIDs(ids)
	if err != nil {
		logger.LogError(err.Error(), "user-app")
		return nil, domain.NotFound("user not found")
	}

	return user, nil
}
