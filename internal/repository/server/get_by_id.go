package ServerRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *ServerRepo) GetById(id int) (*domain.Server, error) {
	var server domain.Server

	query := `SELECT id, owner_id, created_at, name, description FROM servers WHERE id = $1`
	err := r.db.Get(&server, query, id)

	if err != nil {
		return nil, err
	}

	return &server, nil
}
