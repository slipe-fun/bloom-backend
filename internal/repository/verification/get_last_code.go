package verification

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *VerificationRepo) GetLastCode(email string) (*domain.VerificationCode, error) {
	var verificationCode domain.VerificationCode

	query := `SELECT id, email, code, expires_at, created_at FROM verification_codes WHERE email = $1 ORDER BY id DESC`

	start := time.Now()

	err := r.db.Get(&verificationCode, query, email)

	duration := time.Since(start)

	metrics.ObserveDB("verification_code_get_last", duration, err)

	if err != nil {
		return nil, err
	}

	return &verificationCode, nil
}
