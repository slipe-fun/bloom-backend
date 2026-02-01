package ServerRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *ServerRepo) Edit(server *domain.Server) error {
	query := `
		UPDATE servers
		SET name = $1,
			description = $2
			WHERE id = $3
		RETURNING id, owner_id, created_at, name, description
	`

	_, err := r.db.Exec(query, server.Name, server.Description, server.ID)

	if err != nil {
		return err
	}

	return nil
}
