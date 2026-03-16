package chat

import "github.com/lib/pq"

func (r *ChatRepo) GetExistingChatUsers(userID int, ids []int) ([]int, error) {
	query := `
	SELECT DISTINCT (m2->>'id')::int AS user_id
	FROM chats c
	JOIN jsonb_array_elements(c.members) m1 ON TRUE
	JOIN jsonb_array_elements(c.members) m2 ON TRUE
	WHERE (m1->>'id')::int = $1
	AND (m2->>'id')::int = ANY($2::int[])
	AND (m2->>'id')::int != $1;
	`

	rows, err := r.db.Query(query, userID, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []int

	for rows.Next() {
		var id int
		rows.Scan(&id)
		result = append(result, id)
	}

	return result, nil
}
