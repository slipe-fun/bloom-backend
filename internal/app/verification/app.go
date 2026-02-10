package verification

import (
	"github.com/slipe-fun/skid-backend/internal/repository/verification"
)

type VerificationApp struct {
	verification *verification.VerificationRepo
}

func NewAuthApp(verification *verification.VerificationRepo) *VerificationApp {
	return &VerificationApp{
		verification: verification,
	}
}
