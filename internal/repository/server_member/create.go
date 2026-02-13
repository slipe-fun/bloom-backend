package server

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ServerMemberRepo) Create(serverMember *domain.ServerMember) (*domain.ServerMember, error) {
	query := `INSERT INTO server_members (server_id, member_id)
			  VALUES ($1, $2)
			  RETURNING id, server_id, member_id, joined_at`

	var created domain.ServerMember

	start := time.Now()

	err := r.db.QueryRow(query, serverMember.ServerID, serverMember.MemberID).
		Scan(&created.ID, &created.ServerID, &created.MemberID, &created.JoinedAt)

	duration := time.Since(start)

	metrics.ObserveDB("server_member_create", duration, err)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
