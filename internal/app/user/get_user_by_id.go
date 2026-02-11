package user

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (u *UserApp) GetUserByID(id int) (*domain.User, error) {
	user, err := u.users.GetByID(id)
	if err != nil {
		logger.LogError(err.Error(), "user-app")
		return nil, domain.NotFound("user not found")
	}

	return user, nil
}
