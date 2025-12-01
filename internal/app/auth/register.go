package AuthApp

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/service"
)

func (a *AuthApp) Register(email string) error {
	_, err := a.users.GetByEmail(email)

	if err == nil {
		return domain.AlreadyExists("user already exists")
	}

	user, err := a.users.Create(&domain.User{
		Email:       service.Strptr(email),
		DisplayName: service.Strptr(service.GenerateNickname()),
		Username:    "",
	})

	if err != nil {
		return domain.Failed("failed to register user")
	}

	createAndSendCodeError := a.codesApp.CreateAndSendCode(*user.Email)

	if createAndSendCodeError != nil {
		return createAndSendCodeError
	}

	return nil
}
