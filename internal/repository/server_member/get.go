package server

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ServerMemberRepo) Get(server_id, member_id int) (*domain.ServerMember, error) {
	query := `SELECT id, server_id, member_id, joined_at FROM server_members WHERE server_id = $1 AND member_id = $2`

	var serverMember domain.ServerMember

	start := time.Now()

	err := r.db.Get(&serverMember, query, server_id, member_id)

	duration := time.Since(start)

	metrics.ObserveDB("server_member_get", duration, err)

	if err != nil {
		return nil, err
	}

	return &serverMember, nil
}
