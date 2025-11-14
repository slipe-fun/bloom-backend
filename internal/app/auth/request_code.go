package AuthApp

import (
	"errors"
	"time"
)

func (a *AuthApp) RequestCode(email string) error {
	user, err := a.users.GetByEmail(email)

	if err != nil {
		return errors.New("code not found")
	}

	verificationCode, err := a.codesRepo.GetLastCode(email)
	if err == nil {
		now := time.Now().UTC()

		if !verificationCode.ExpiresAt.Before(now) {
			return errors.New("code hasnt expired")
		}

		err = a.codesRepo.DeleteByEmailAndCode(email, verificationCode.Code)
		if err != nil {
			return err
		}
	}

	createAndSendCodeError := a.codesApp.CreateAndSendCode(user.Email)

	if createAndSendCodeError != nil {
		return createAndSendCodeError
	}

	return nil
}
