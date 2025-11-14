package AuthApp

import (
	"errors"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (a *AuthApp) Register(username string, email string) error {
	_, err := a.users.GetByUsername(username)

	if err == nil {
		return errors.New("user already exists")
	}

	_, err = a.users.GetByEmail(email)

	if err == nil {
		return errors.New("user already exists")
	}

	user, err := a.users.Create(&domain.User{
		Username: username,
		Email:    email,
	})

	if err != nil {
		return err
	}

	createAndSendCodeError := a.codesApp.CreateAndSendCode(user.Email)

	if createAndSendCodeError != nil {
		return createAndSendCodeError
	}

	return nil
}
