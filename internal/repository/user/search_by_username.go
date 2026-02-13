package user

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *UserRepo) SearchUsersByUsername(query string, limit, offset int) ([]*domain.User, error) {
	sqlQuery := `
	SELECT id, username, display_name, description, date
	FROM users
	WHERE 
		username ILIKE '%' || $1 || '%'
		OR display_name ILIKE '%' || $1 || '%'
		OR similarity(username, cyr_to_lat($1)) > 0.3
		OR similarity(cyr_to_lat(display_name), cyr_to_lat($1)) > 0.3
	ORDER BY
		CASE WHEN username ILIKE $1 THEN 1 ELSE 0 END DESC,
		
		CASE WHEN username ILIKE $1 || '%' THEN 1 ELSE 0 END DESC,
		
		CASE WHEN display_name ILIKE $1 || '%' THEN 1 ELSE 0 END DESC,

		GREATEST(
			similarity(username, cyr_to_lat($1)), 
			similarity(cyr_to_lat(display_name), cyr_to_lat($1))
		) DESC
	LIMIT $2 OFFSET $3;
	`

	start := time.Now()

	rows, err := r.db.Query(sqlQuery, query, limit, offset)

	duration := time.Since(start)

	metrics.ObserveDB("users_search_by_username_and_displayname", duration, err)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.DisplayName, &user.Description, &user.Date); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
