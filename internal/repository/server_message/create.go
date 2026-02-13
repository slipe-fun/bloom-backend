package server

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ServerMessageRepo) Create(serverMessage *domain.ServerMessage) (*domain.ServerMessage, error) {
	query := `INSERT INTO server_messages (server_id, channel_id, member_id, content)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, server_id, channel_id, member_id, content, sent_at`

	var created domain.ServerMessage

	start := time.Now()

	err := r.db.QueryRow(query, serverMessage.ServerID, serverMessage.ChannelID, serverMessage.MemberID, serverMessage.Content).
		Scan(&created.ID, &created.ServerID, &created.ChannelID, &created.MemberID, &created.Content, &created.SentAt)

	duration := time.Since(start)

	metrics.ObserveDB("server_message_create", duration, err)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
