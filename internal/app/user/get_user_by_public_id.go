package user

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (u *UserApp) GetUserByPublicID(id string) (*domain.User, error) {
	user, err := u.users.GetByPublicID(id)
	if err != nil {
		logger.LogError(err.Error(), "user-app")
		return nil, domain.NotFound("user not found")
	}

	if user.KyberPublicKey == "" || user.EcdhPublicKey == "" || user.EdPublicKey == "" {
		return nil, domain.NotFound("user not found")
	}

	return user, nil
}
