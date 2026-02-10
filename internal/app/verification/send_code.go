package verification

import (
	"fmt"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/generator"
	"github.com/slipe-fun/skid-backend/internal/mailer"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (v *VerificationApp) CreateAndSendCode(email string) error {
	code, err := generator.GenerateNumericCode(6)
	if err != nil {
		logger.LogError(err.Error(), "verification-app")
		return domain.Failed("failed to generate numeric code")
	}

	createdCode, err := v.verification.Create(&domain.VerificationCode{
		Email: email,
		Code:  code,
	})
	if err != nil {
		logger.LogError(err.Error(), "verification-app")
		return domain.Failed("failed to create code")
	}

	sendEmailError := mailer.SendMail(
		fmt.Sprintf("Your code - %s", createdCode.Code),
		fmt.Sprintf("Hello! Your Bloom verification code - %s", createdCode.Code),
		email,
	)
	if sendEmailError != nil {
		logger.LogError(sendEmailError.Error(), "verification-app")
		v.verification.DeleteByEmailAndCode(email, createdCode.Code)
		return domain.Failed("failed to send email")
	}

	return nil
}
