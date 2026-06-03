package user

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *UserRepo) Edit(user *domain.User) error {
	query := `
        UPDATE users
        SET username = $1,
			display_name = $2,
			description = $3
        WHERE id = $4
        RETURNING id, public_id, username, display_name, description, date
    `

	start := time.Now()

	_, err := r.db.Exec(query, user.Username, user.DisplayName, user.Description, user.ID)

	duration := time.Since(start)

	metrics.ObserveDB("user_edit", duration, err)

	if err != nil {
		return err
	}

	return nil
}
