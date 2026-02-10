package verification

import "github.com/slipe-fun/skid-backend/internal/domain"

type VerificationRepo interface {
	Create(code *domain.VerificationCode) (*domain.VerificationCode, error)
	DeleteByEmailAndCode(email, code string) error
}
