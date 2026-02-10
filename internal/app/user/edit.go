package user

import (
	"database/sql"
	"errors"

	"github.com/slipe-fun/skid-backend/internal/config"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (u *UserApp) EditUser(token string, editedUser *domain.User) (*domain.User, error) {
	session, err := u.sessionApp.GetSession(token)
	if err != nil {
		return nil, err
	}

	user, err := u.users.GetById(editedUser.ID)
	if err != nil {
		logger.LogError(err.Error(), "user-app")
		return nil, domain.Failed("failed to get user")
	}

	if session.UserID != user.ID {
		return nil, domain.Failed("failed to get user")
	}

	if editedUser.Username == "" {
		user.Username = ""
	} else {
		if !config.UsernameRegex.MatchString(editedUser.Username) {
			return nil, domain.InvalidData("invalid username")
		}

		existing, err := u.users.GetByUsername(editedUser.Username)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.LogError(err.Error(), "user-app")
			return nil, err
		}

		if existing != nil && existing.ID != editedUser.ID {
			logger.LogError(err.Error(), "user-app")
			return nil, domain.AlreadyExists("user with this username already exists")
		}

		user.Username = editedUser.Username
	}

	if editedUser.DisplayName != nil {
		dLen := len(*editedUser.DisplayName)
		if dLen > 0 && dLen <= 20 {
			user.DisplayName = editedUser.DisplayName
		} else {
			return nil, domain.InvalidData("invalid length of display name")
		}
	}

	if editedUser.Description != nil {
		if len(*editedUser.Description) <= 150 {
			user.Description = editedUser.Description
		} else {
			return nil, domain.InvalidData("invalid length of description")
		}
	}

	editUserErr := u.users.Edit(user)
	if editUserErr != nil {
		logger.LogError(editUserErr.Error(), "user-app")
		return nil, domain.Failed("failed to edit user")
	}

	return user, nil
}
