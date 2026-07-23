package user

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *UserRepo) Create(user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (public_id, auth_lookup_id, username, display_name, description, ml_kem_public_key, ecdh_public_key, ed_public_key)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	          RETURNING id, public_id, auth_lookup_id, username, display_name, description, ml_kem_public_key, ecdh_public_key, ed_public_key, date`

	var created domain.User

	start := time.Now()

	err := r.db.QueryRow(query, user.PublicID, user.AuthLookupID, user.Username, user.DisplayName, user.Description, user.MlKemPublicKey, user.EcdhPublicKey, user.EdPublicKey).
		Scan(&created.ID, &created.PublicID, &user.AuthLookupID, &created.Username, &created.DisplayName, &created.Description, &created.MlKemPublicKey, &created.EcdhPublicKey, &created.EdPublicKey, &created.Date)

	duration := time.Since(start)

	metrics.ObserveDB("user_create", duration, err)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
