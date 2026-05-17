package credential

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *CredentialRepo) UpdateSignCount(credentialID []byte, signCount uint32, cloneWarning bool) error {
	query := `UPDATE user_credentials SET sign_count = $1, clone_warning = $2 WHERE credential_id = $3`

	start := time.Now()

	_, err := r.db.Exec(query, signCount, cloneWarning, credentialID)

	duration := time.Since(start)

	metrics.ObserveDB("credential_update_sign_count", duration, err)

	return err
}
