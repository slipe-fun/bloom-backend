package AuthApp

import (
	"errors"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/service"
)

func (a *AuthApp) Register(email string) error {
	_, err := a.users.GetByEmail(email)

	if err == nil {
		return errors.New("user already exists")
	}

	user, err := a.users.Create(&domain.User{
		Email:       email,
		DisplayName: service.GenerateNickname(),
		Username:    "",
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
