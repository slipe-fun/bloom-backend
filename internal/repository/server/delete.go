package server

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ServerRepo) Delete(id int) error {
	query := `DELETE FROM servers WHERE id = $1`

	start := time.Now()

	_, err := r.db.Exec(query, id)

	duration := time.Since(start)

	metrics.ObserveDB("server_delete", duration, err)

	if err != nil {
		return err
	}

	return nil
}
