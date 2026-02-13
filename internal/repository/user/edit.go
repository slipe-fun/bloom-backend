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
            email = $2,
			display_name = $3,
			description = $4
        WHERE id = $5
        RETURNING id, username, email, display_name, description, date
    `

	start := time.Now()

	_, err := r.db.Exec(query, user.Username, user.Email, user.DisplayName, user.Description, user.ID)

	duration := time.Since(start)

	metrics.ObserveDB("user_edit", duration, err)

	if err != nil {
		return err
	}

	return nil
}
