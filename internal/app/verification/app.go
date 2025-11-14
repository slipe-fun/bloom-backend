package VerificationApp

import (
	VerificationRepo "github.com/slipe-fun/skid-backend/internal/repository/verification"
)

type VerificationApp struct {
	verification *VerificationRepo.VerificationRepo
}

func NewAuthApp(verification *VerificationRepo.VerificationRepo) *VerificationApp {
	return &VerificationApp{
		verification: verification,
	}
}
