package user

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *UserRepo) GetByAuthLookupID(authLookupID string) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, public_id, auth_lookup_id, username, display_name, description, ml_kem_public_key, ecdh_public_key, ed_public_key, date FROM users WHERE auth_lookup_id = $1`

	start := time.Now()

	err := r.db.Get(&user, query, authLookupID)

	duration := time.Since(start)

	metrics.ObserveDB("user_get_by_id", duration, err)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
