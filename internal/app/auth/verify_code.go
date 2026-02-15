package auth

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (a *AuthApp) VerifyCode(email string, code string) (string, *domain.Session, *domain.User, error) {
	user, err := a.users.GetByEmail(email)

	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", nil, nil, domain.NotFound("code not found")
	}

	verificationCode, err := a.codesRepo.GetByEmailAndCode(email, code)
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", nil, nil, domain.NotFound("code not found")
	}

	now := time.Now().UTC()

	if verificationCode.ExpiresAt.Before(now) {
		return "", nil, nil, domain.Expired("code has expired")
	}

	token, session, err := a.sessionApp.CreateSession(user.ID)
	if err != nil {
		return "", nil, nil, err
	}

	err = a.codesRepo.DeleteByEmailAndCode(email, code)
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", nil, nil, domain.Failed("failed to delete old code")
	}

	return token, session, user, nil
}
