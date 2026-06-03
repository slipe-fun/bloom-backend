package user

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *UserRepo) GetAllUsers(limit, offset int) ([]*domain.User, error) {
	query := `SELECT id, public_id, username, display_name, description, kyber_public_key, ecdh_public_key, ed_public_key, date
			  FROM users
			  WHERE kyber_public_key IS NOT NULL AND kyber_public_key != ''
					AND ecdh_public_key IS NOT NULL AND ecdh_public_key != ''
					AND ed_public_key IS NOT NULL AND ed_public_key != ''
			  ORDER BY id DESC
			  LIMIT $1 OFFSET $2`

	start := time.Now()

	rows, err := r.db.Query(query, limit, offset)

	duration := time.Since(start)

	metrics.ObserveDB("users_get_all", duration, err)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.PublicID, &user.Username, &user.DisplayName, &user.Description, &user.KyberPublicKey, &user.EcdhPublicKey, &user.EdPublicKey, &user.Date); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
