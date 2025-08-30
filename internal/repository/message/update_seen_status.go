package MessageRepo

import (
	"fmt"
	"time"
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

	_, err := r.db.Exec(query, args...)
	return err
}
