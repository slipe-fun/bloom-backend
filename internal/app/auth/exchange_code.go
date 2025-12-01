package AuthApp

import (
	"database/sql"
	"errors"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/service"
)

func (a *AuthApp) ExchangeCode(code string) (string, *domain.User, error) {
	token, err := a.google.ExchangeCode(code)
	if err != nil {
		return "", nil, domain.Failed("failed to exchange code")
	}

	client, err := a.google.GetUserInfo(token)
	if err != nil {
		return "", nil, domain.Failed("failed to get user info")
	}

	email, ok := client["email"].(string)
	if !ok {
		return "", nil, domain.InvalidData("email not found or wrong type")
	}

	user, err := a.users.GetByEmail(email)
	if errors.Is(err, sql.ErrNoRows) {
		name, ok := client["name"].(string)
		if !ok {
			return "", nil, domain.InvalidData("name not found or wrong type")
		}

		user, err = a.users.Create(&domain.User{
			Email:       service.Strptr(email),
			Username:    service.GenerateUsername(name),
			DisplayName: service.Strptr(service.GenerateNickname()),
		})
		if err != nil {
			return "", nil, domain.Failed("failed to register user")
		}
	} else if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return "", nil, domain.Failed("failed to get user")
	}

	jwtToken, err := a.sessionApp.CreateSession(user.ID)
	if err != nil {
		return "", nil, domain.Failed("failed to generate token")
	}

	return jwtToken, user, nil
}
