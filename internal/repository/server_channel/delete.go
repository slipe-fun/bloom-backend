package server

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ServerChannelRepo) Delete(id int) error {
	query := `DELETE FROM server_channels WHERE id = $1`

	start := time.Now()

	_, err := r.db.Exec(query, id)

	duration := time.Since(start)

	metrics.ObserveDB("server_channel_delete", duration, err)

	if err != nil {
		return err
	}

	return nil
}
