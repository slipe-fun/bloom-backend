package server

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ServerRepo) Create(server *domain.Server) (*domain.Server, error) {
	query := `INSERT INTO servers (owner_id, name, description)
			  VALUES ($1, $2, $3)
			  RETURNING id, owner_id, created_at, name, description`

	var created domain.Server

	start := time.Now()

	err := r.db.QueryRow(query, server.OwnerID, server.Name, server.Description).
		Scan(&created.ID, &created.OwnerID, &created.CreatedAt, &created.Name, &created.Description)

	duration := time.Since(start)

	metrics.ObserveDB("server_create", duration, err)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
