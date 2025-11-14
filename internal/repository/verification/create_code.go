package VerificationRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *VerificationRepo) Create(code *domain.VerificationCode) (*domain.VerificationCode, error) {
	query := `INSERT INTO verification_codes (email, code) 
	          VALUES ($1, $2) 
	          RETURNING id, email, code, expires_at, created_at`

	var created domain.VerificationCode
	err := r.db.QueryRow(query, code.Email, code.Code).
		Scan(&created.ID, &created.Email, &created.Code, &created.ExpiresAt, &created.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
