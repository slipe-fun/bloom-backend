package ServerMemberRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *ServerMemberRepo) Create(serverMember *domain.ServerMember) (*domain.ServerMember, error) {
	query := `INSERT INTO server_members (server_id, member_id)
			  VALUES ($1, $2)
			  RETURNING id, server_id, member_id, joined_at`

	var created domain.ServerMember
	err := r.db.QueryRow(query, serverMember.ServerID, serverMember.MemberID).
		Scan(&created.ID, &created.ServerID, &created.MemberID, &created.JoinedAt)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
