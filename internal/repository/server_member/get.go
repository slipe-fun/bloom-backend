package server

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *ServerMemberRepo) Get(server_id, member_id int) (*domain.ServerMember, error) {
	query := `SELECT id, server_id, member_id, joined_at FROM server_members WHERE server_id = $1 AND member_id = $2`

	var serverMember domain.ServerMember
	err := r.db.Get(&serverMember, query, server_id, member_id)

	if err != nil {
		return nil, err
	}

	return &serverMember, nil
}
