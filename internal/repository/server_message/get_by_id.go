package server

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ServerMessageRepo) GetByID(id int) (*domain.ServerMessage, error) {
	query := `SELECT id, server_id, member_id, channel_id, content, sent_at FROM server_messages WHERE id = $1`

	start := time.Now()

	row := r.db.QueryRow(query, id)

	var message domain.ServerMessage
	err := row.Scan(&message.ID, &message.ServerID, &message.MemberID, &message.ChannelID, &message.Content, &message.SentAt)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start)

	metrics.ObserveDB("server_message_get_by_id", duration, err)

	return &message, nil
}
