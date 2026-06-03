package user

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *UserRepo) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	start := time.Now()

	_, err := r.db.Exec(query, id)

	duration := time.Since(start)

	metrics.ObserveDB("user_delete", duration, err)

	if err != nil {
		return err
	}

	return nil
}
