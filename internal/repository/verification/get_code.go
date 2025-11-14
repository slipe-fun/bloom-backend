package VerificationRepo

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *VerificationRepo) GetByEmailAndCode(email string, code string) (*domain.VerificationCode, error) {
	var verificationCode domain.VerificationCode

	query := `SELECT id, email, code, expires_at, created_at FROM verification_codes WHERE email = $1 AND code = $2`
	err := r.db.Get(&verificationCode, query, email, code)

	if err != nil {
		return nil, err
	}

	return &verificationCode, nil
}
