package ServerChannelRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *ServerChannelRepo) Edit(ServerChannel *domain.ServerChannel) error {
	query := `UPDATE server_channels
		SET name = $1,
			position = $2
		WHERE id = $3
		RETURNING id, server_id, name, type, position, created_at
	`

	_, err := r.db.Exec(query, ServerChannel.Name, ServerChannel.Position, ServerChannel.ID)

	if err != nil {
		return err
	}

	return nil
}
