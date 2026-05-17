package credential

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *CredentialRepo) GetByCredentialID(credentialID []byte) (*domain.Credential, error) {
	var credential domain.Credential

	query := `SELECT id, user_id, credential_id, public_key, attestation_type, sign_count, clone_warning, transport, created_at 
	          FROM user_credentials WHERE credential_id = $1`

	start := time.Now()

	err := r.db.Get(&credential, query, credentialID)

	duration := time.Since(start)

	metrics.ObserveDB("credential_get_by_credential_id", duration, err)

	if err != nil {
		return nil, err
	}

	return &credential, nil
}
