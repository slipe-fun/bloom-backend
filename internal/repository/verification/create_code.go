package verification

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *VerificationRepo) Create(code *domain.VerificationCode) (*domain.VerificationCode, error) {
	query := `INSERT INTO verification_codes (email, code) 
	          VALUES ($1, $2) 
	          RETURNING id, email, code, expires_at, created_at`

	var created domain.VerificationCode

	start := time.Now()

	err := r.db.QueryRow(query, code.Email, code.Code).
		Scan(&created.ID, &created.Email, &created.Code, &created.ExpiresAt, &created.CreatedAt)

	duration := time.Since(start)

	metrics.ObserveDB("verification_code_create", duration, err)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
