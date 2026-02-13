package server

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ServerRepo) Edit(server *domain.Server) error {
	query := `
		UPDATE servers
		SET name = $1,
			description = $2
			WHERE id = $3
		RETURNING id, owner_id, created_at, name, description
	`

	start := time.Now()

	_, err := r.db.Exec(query, server.Name, server.Description, server.ID)

	duration := time.Since(start)

	metrics.ObserveDB("server_edit", duration, err)

	if err != nil {
		return err
	}

	return nil
}
