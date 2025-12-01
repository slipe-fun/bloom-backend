package UserApp

import "github.com/slipe-fun/skid-backend/internal/domain"

func (u *UserApp) GetUserById(id int) (*domain.User, error) {
	user, err := u.users.GetById(id)
	if err != nil {
		return nil, domain.NotFound("user not found")
	}

	return user, nil
}
