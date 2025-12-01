package AuthApp

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (a *AuthApp) VerifyCode(email string, code string) (string, *domain.User, error) {
	user, err := a.users.GetByEmail(email)

	if err != nil {
		return "", nil, domain.NotFound("code not found")
	}

	verificationCode, err := a.codesRepo.GetByEmailAndCode(email, code)
	if err != nil {
		return "", nil, domain.NotFound("code not found")
	}

	now := time.Now().UTC()

	if verificationCode.ExpiresAt.Before(now) {
		return "", nil, domain.Expired("code has expired")
	}

	token, err := a.sessionApp.CreateSession(user.ID)
	if err != nil {
		return "", nil, err
	}

	err = a.codesRepo.DeleteByEmailAndCode(email, code)
	if err != nil {
		return "", nil, domain.Failed("failed to delete old code")
	}

	return token, user, nil
}
