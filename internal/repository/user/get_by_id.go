package user

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *UserRepo) GetByID(id int) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, username, email, display_name, description, date FROM users WHERE id = $1`

	start := time.Now()

	err := r.db.Get(&user, query, id)

	duration := time.Since(start)

	metrics.ObserveDB("user_get_by_id", duration, err)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
