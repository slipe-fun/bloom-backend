package server

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *ServerChannelRepo) GetByID(id int) (*domain.ServerChannel, error) {
	query := `SELECT id, server_id, name, type, position, created_at FROM server_channels WHERE id = $1`

	var serverChannel domain.ServerChannel
	err := r.db.Get(&serverChannel, query, id)

	if err != nil {
		return nil, err
	}

	return &serverChannel, nil
}
