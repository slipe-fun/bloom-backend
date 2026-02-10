package user

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (u *UserApp) IsUserWithEmailExists(email string) (bool, error) {
	_, err := u.users.GetByEmail(email)
	if err != nil {
		logger.LogError(err.Error(), "user-app")
		return false, domain.Failed("user not found")
	}

	return true, nil
}
