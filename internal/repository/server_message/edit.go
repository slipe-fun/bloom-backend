package server

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ServerMessageRepo) Edit(id int, newContent string) error {
	query := `UPDATE server_messages SET content = $1 WHERE id = $2`

	start := time.Now()

	_, err := r.db.Exec(query, newContent, id)

	duration := time.Since(start)

	metrics.ObserveDB("server_message_edit", duration, err)

	if err != nil {
		return err
	}

	return nil
}
