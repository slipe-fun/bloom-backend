package server

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ServerRepo) GetByID(id int) (*domain.Server, error) {
	var server domain.Server

	query := `SELECT id, owner_id, created_at, name, description FROM servers WHERE id = $1`

	start := time.Now()

	err := r.db.Get(&server, query, id)

	duration := time.Since(start)

	metrics.ObserveDB("server_get_by_id", duration, err)

	if err != nil {
		return nil, err
	}

	return &server, nil
}
