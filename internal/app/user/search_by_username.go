package UserApp

import "github.com/slipe-fun/skid-backend/internal/domain"

func (u *UserApp) SearchUsersByUsername(username string, limit, offset int) ([]*domain.User, error) {
	users, err := u.users.SearchUsersByUsername(username, limit, offset)
	if err != nil {
		return nil, domain.NotFound("users not found")
	}

	if len(users) == 0 {
		return nil, domain.NotFound("users not found")
	}

	return users, nil
}
