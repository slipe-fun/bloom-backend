package user

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *UserRepo) Create(user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (username, display_name, description) 
	          VALUES ($1, $2, $3) 
	          RETURNING id, username, display_name, description, date`

	var created domain.User

	start := time.Now()

	err := r.db.QueryRow(query, user.Username, user.DisplayName, user.Description).
		Scan(&created.ID, &created.Username, &created.DisplayName, &created.Description, &created.Date)

	duration := time.Since(start)

	metrics.ObserveDB("user_create", duration, err)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
