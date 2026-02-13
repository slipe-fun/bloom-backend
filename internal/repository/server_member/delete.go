package server

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ServerMemberRepo) Delete(server_id, member_id int) error {
	query := `DELETE FROM server_members WHERE server_id = $1 AND member_id = $2`

	start := time.Now()

	_, err := r.db.Exec(query, server_id, member_id)

	duration := time.Since(start)

	metrics.ObserveDB("server_member_delete", duration, err)

	if err != nil {
		return err
	}

	return nil
}
