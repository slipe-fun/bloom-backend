package user

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *UserRepo) Create(user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (public_id, username, display_name, description, kyber_public_key, ecdh_public_key, ed_public_key) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) 
	          RETURNING id, public_id, username, display_name, description, kyber_public_key, ecdh_public_key, ed_public_key, date`

	var created domain.User

	start := time.Now()

	err := r.db.QueryRow(query, user.PublicID, user.Username, user.DisplayName, user.Description, user.KyberPublicKey, user.EcdhPublicKey, user.EdPublicKey).
		Scan(&created.ID, &created.PublicID, &created.Username, &created.DisplayName, &created.Description, &created.KyberPublicKey, &created.EcdhPublicKey, &created.EdPublicKey, &created.Date)

	duration := time.Since(start)

	metrics.ObserveDB("user_create", duration, err)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
