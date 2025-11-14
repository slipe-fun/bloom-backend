package AuthApp

import (
	"errors"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (a *AuthApp) VerifyCode(email string, code string) (string, *domain.User, error) {
	user, err := a.users.GetByEmail(email)

	if err != nil {
		return "", nil, errors.New("code not found")
	}

	verificationCode, err := a.codesRepo.GetByEmailAndCode(email, code)
	if err != nil {
		return "", nil, errors.New("code not found")
	}

	now := time.Now().UTC()

	if verificationCode.ExpiresAt.Before(now) {
		return "", nil, errors.New("code has expired")
	}

	token, err := a.jwtSvc.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	err = a.codesRepo.DeleteByEmailAndCode(email, code)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
