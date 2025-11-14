package VerificationRepo

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *VerificationRepo) GetLastCode(email string) (*domain.VerificationCode, error) {
	var verificationCode domain.VerificationCode

	query := `SELECT id, email, code, expires_at, created_at FROM verification_codes WHERE email = $1 ORDER BY id DESC`
	err := r.db.Get(&verificationCode, query, email)

	if err != nil {
		return nil, err
	}

	return &verificationCode, nil
}
