package ServerChannelRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *ServerChannelRepo) Create(serverChannel *domain.ServerChannel) (*domain.ServerChannel, error) {
	query := `INSERT INTO server_channels (server_id, name, type, position)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, server_id, name, type, position, created_at`

	var created domain.ServerChannel
	err := r.db.QueryRow(query, serverChannel.ServerID, serverChannel.Name, serverChannel.Type, serverChannel.Position).
		Scan(&created.ID, &created.ServerID, &created.Name, &created.Type, &created.Position, &created.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
