package server

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ServerChannelRepo) Edit(ServerChannel *domain.ServerChannel) error {
	query := `UPDATE server_channels
		SET name = $1,
			position = $2
		WHERE id = $3
		RETURNING id, server_id, name, type, position, created_at
	`

	start := time.Now()

	_, err := r.db.Exec(query, ServerChannel.Name, ServerChannel.Position, ServerChannel.ID)

	duration := time.Since(start)

	metrics.ObserveDB("server_channel_edit", duration, err)

	if err != nil {
		return err
	}

	return nil
}
