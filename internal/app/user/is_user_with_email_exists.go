package UserApp

import "github.com/slipe-fun/skid-backend/internal/domain"

func (u *UserApp) IsUserWithEmailExists(email string) (bool, error) {
	_, err := u.users.GetByEmail(email)
	if err != nil {
		return false, domain.Failed("user not found")
	}

	return true, nil
}
