package user

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *UserRepo) GetByPublicID(id string) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, public_id, username, display_name, description, ml_kem_public_key, ecdh_public_key, ed_public_key, date FROM users WHERE public_id = $1`

	start := time.Now()

	err := r.db.Get(&user, query, id)

	duration := time.Since(start)

	metrics.ObserveDB("user_get_by_public_id", duration, err)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
