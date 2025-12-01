package VerificationApp

import (
	"fmt"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/service"
)

func (v *VerificationApp) CreateAndSendCode(email string) error {
	code, err := service.GenerateNumericCode(6)
	if err != nil {
		return domain.Failed("failed to generate numeric code")
	}

	createdCode, err := v.verification.Create(&domain.VerificationCode{
		Email: email,
		Code:  code,
	})
	if err != nil {
		return domain.Failed("failed to create code")
	}

	sendEmailError := service.SendMail(
		fmt.Sprintf("Your code - %s", createdCode.Code),
		fmt.Sprintf("Hello! Your Bloom verification code - %s", createdCode.Code),
		email,
	)
	if sendEmailError != nil {
		v.verification.DeleteByEmailAndCode(email, createdCode.Code)
		return domain.Failed("failed to send email")
	}

	return nil
}
