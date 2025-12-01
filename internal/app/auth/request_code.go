package AuthApp

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (a *AuthApp) RequestCode(email string) error {
	user, err := a.users.GetByEmail(email)

	if err != nil {
		return domain.NotFound("code not found")
	}

	verificationCode, err := a.codesRepo.GetLastCode(email)
	if err == nil {
		now := time.Now().UTC()

		if !verificationCode.ExpiresAt.Before(now) {
			return domain.NotExpired("code hasnt expired")
		}

		err = a.codesRepo.DeleteByEmailAndCode(email, verificationCode.Code)
		if err != nil {
			return domain.Failed("failed to delete old code")
		}
	}

	createAndSendCodeError := a.codesApp.CreateAndSendCode(*user.Email)

	if createAndSendCodeError != nil {
		return createAndSendCodeError
	}

	return nil
}
