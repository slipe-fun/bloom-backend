package server

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ServerChannelRepo) GetByID(id int) (*domain.ServerChannel, error) {
	query := `SELECT id, server_id, name, type, position, created_at FROM server_channels WHERE id = $1`

	var serverChannel domain.ServerChannel

	start := time.Now()

	err := r.db.Get(&serverChannel, query, id)

	duration := time.Since(start)

	metrics.ObserveDB("server_channel_get_by_id", duration, err)

	if err != nil {
		return nil, err
	}

	return &serverChannel, nil
}
