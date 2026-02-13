package verification

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *VerificationRepo) GetByEmailAndCode(email string, code string) (*domain.VerificationCode, error) {
	var verificationCode domain.VerificationCode

	query := `SELECT id, email, code, expires_at, created_at FROM verification_codes WHERE email = $1 AND code = $2`

	start := time.Now()

	err := r.db.Get(&verificationCode, query, email, code)

	duration := time.Since(start)

	metrics.ObserveDB("verification_code_get_by_email_and_code", duration, err)

	if err != nil {
		return nil, err
	}

	return &verificationCode, nil
}
