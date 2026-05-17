package credential

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *CredentialRepo) Create(credential *domain.Credential) (*domain.Credential, error) {
	query := `INSERT INTO user_credentials (user_id, credential_id, public_key, attestation_type, sign_count, clone_warning, transport) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) 
	          RETURNING id, user_id, credential_id, public_key, attestation_type, sign_count, clone_warning, transport, created_at`

	var created domain.Credential

	start := time.Now()

	err := r.db.QueryRow(query, credential.UserID, credential.CredentialID, credential.PublicKey, credential.AttestationType, credential.SignCount, credential.CloneWarning, credential.Transport).
		Scan(&created.ID, &created.UserID, &created.CredentialID, &created.PublicKey, &created.AttestationType, &created.SignCount, &created.CloneWarning, &created.Transport, &created.CreatedAt)

	duration := time.Since(start)

	metrics.ObserveDB("credential_create", duration, err)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
