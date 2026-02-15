package auth

import (
	"database/sql"
	"errors"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/generator"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
	"github.com/slipe-fun/skid-backend/internal/pointer"
)

func (a *AuthApp) ExchangeCode(code string) (string, *domain.Session, *domain.User, error) {
	token, err := a.google.ExchangeCode(code)
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", nil, nil, domain.Failed("failed to exchange code")
	}

	client, err := a.google.GetUserInfo(token)
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", nil, nil, domain.Failed("failed to get user info")
	}

	email, ok := client["email"].(string)
	if !ok {
		return "", nil, nil, domain.InvalidData("email not found or wrong type")
	}

	user, err := a.users.GetByEmail(email)
	if errors.Is(err, sql.ErrNoRows) {
		name, ok := client["name"].(string)
		if !ok {
			return "", nil, nil, domain.InvalidData("name not found or wrong type")
		}

		user, err = a.users.Create(&domain.User{
			Email:       pointer.Strptr(email),
			Username:    generator.GenerateUsername(name),
			DisplayName: pointer.Strptr(generator.GenerateNickname()),
		})
		if err != nil {
			logger.LogError(err.Error(), "auth-app")
			return "", nil, nil, domain.Failed("failed to register user")
		}
	} else if !errors.Is(err, sql.ErrNoRows) && err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", nil, nil, domain.Failed("failed to get user")
	}

	jwtToken, session, err := a.sessionApp.CreateSession(user.ID)
	if err != nil {
		return "", nil, nil, domain.Failed("failed to generate token")
	}

	return jwtToken, session, user, nil
}
