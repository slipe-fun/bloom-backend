package auth

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/generator"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
	"github.com/slipe-fun/skid-backend/internal/pointer"
)

func (a *AuthApp) Register(email string) error {
	_, err := a.users.GetByEmail(email)

	if err == nil {
		return domain.AlreadyExists("user already exists")
	}

	user, err := a.users.Create(&domain.User{
		Email:       pointer.Strptr(email),
		DisplayName: pointer.Strptr(generator.GenerateNickname()),
		Username:    "",
	})

	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return domain.Failed("failed to register user")
	}

	createAndSendCodeError := a.codesApp.CreateAndSendCode(*user.Email)

	if createAndSendCodeError != nil {
		return createAndSendCodeError
	}

	return nil
}
