package UserApp

import (
	"database/sql"
	"errors"

	"github.com/slipe-fun/skid-backend/internal/config"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (u *UserApp) EditUser(token string, editedUser *domain.User) (*domain.User, error) {
	session, err := u.sessionApp.GetSession(token)
	if err != nil {
		return nil, err
	}

	user, err := u.users.GetById(editedUser.ID)
	if err != nil {
		return nil, errors.New("failed to get user")
	}

	if session.UserID != user.ID {
		return nil, errors.New("failed to get user")
	}

	if editedUser.Username == "" {
		user.Username = ""
	} else {
		if !config.UsernameRegex.MatchString(editedUser.Username) {
			return nil, errors.New("wrong username")
		}

		existing, err := u.users.GetByUsername(editedUser.Username)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		if existing != nil && existing.ID != editedUser.ID {
			return nil, errors.New("user with this username already exists")
		}

		user.Username = editedUser.Username
	}

	if len(editedUser.DisplayName) > 0 && len(editedUser.DisplayName) <= 20 {
		user.DisplayName = editedUser.DisplayName
	} else {
		return nil, errors.New("wrong length of display name")
	}

	editUserErr := u.users.Edit(user)
	if editUserErr != nil {
		return nil, errors.New("failed to edit user")
	}

	return user, nil
}
