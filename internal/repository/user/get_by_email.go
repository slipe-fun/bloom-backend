package user

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *UserRepo) GetByEmail(email string) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, username, email, display_name, description, date FROM users WHERE email = $1`

	start := time.Now()

	err := r.db.Get(&user, query, email)

	duration := time.Since(start)

	metrics.ObserveDB("create_friend", duration, err)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
