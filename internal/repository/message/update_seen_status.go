package message

import (
	"fmt"
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *MessageRepo) UpdateMessagesSeenStatus(messages []int, seenTime time.Time) error {
	if len(messages) == 0 {
		return nil
	}

	placeholders := ""
	args := []interface{}{seenTime}
	for i, id := range messages {
		if i > 0 {
			placeholders += ", "
		}
		args = append(args, id)
		placeholders += fmt.Sprintf("$%d", i+2)
	}

	query := fmt.Sprintf(`
		UPDATE messages
		SET seen = $1
		WHERE id IN (%s)
	`, placeholders)

	start := time.Now()

	_, err := r.db.Exec(query, args...)

	duration := time.Since(start)

	metrics.ObserveDB("message_update_seen_status", duration, err)

	return err
}
