package user

import (
	"time"

	"github.com/lib/pq"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *UserRepo) GetUsersByIDs(ids []int) ([]domain.User, error) {
	var users []domain.User

	query := `
		SELECT
			id,
			public_id,
			auth_lookup_id,
			username,
			display_name,
			description,
			ml_kem_public_key,
			ecdh_public_key,
			ed_public_key,
			date
		FROM users
		WHERE id = ANY($1)
	`

	start := time.Now()

	err := r.db.Select(&users, query, pq.Array(ids))

	duration := time.Since(start)

	metrics.ObserveDB("user_get_by_ids", duration, err)

	if err != nil {
		return nil, err
	}

	return users, nil
}
