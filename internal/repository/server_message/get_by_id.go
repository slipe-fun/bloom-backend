package ServerMessageRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *ServerMessageRepo) GetById(id int) (*domain.ServerMessage, error) {
	query := `SELECT id, server_id, member_id, channel_id, content, sent_at FROM server_messages WHERE id = $1`

	row := r.db.QueryRow(query, id)

	var message domain.ServerMessage
	err := row.Scan(&message.ID, &message.ServerID, &message.MemberID, &message.ChannelID, &message.Content, &message.SentAt)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
