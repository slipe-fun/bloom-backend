package session

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *SessionRepo) Delete(id int) error {
	query := `DELETE FROM sessions WHERE id = $1`

	start := time.Now()

	_, err := r.db.Exec(query, id)

	duration := time.Since(start)

	metrics.ObserveDB("session_delete", duration, err)

	if err != nil {
		return err
	}

	return nil
}
