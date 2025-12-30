package UserRepo

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *UserRepo) SearchUsersByUsername(username string, limit, offset int) ([]*domain.User, error) {
	rows, err := r.db.Query(`
	SELECT id, username, display_name, description, date
	FROM users
	WHERE similarity(cyr_to_lat(username), cyr_to_lat($1)) > 0.3
	ORDER BY similarity(cyr_to_lat(username), cyr_to_lat($1)) DESC
	LIMIT $2 OFFSET $3;
	`, username, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.DisplayName, &user.Description, &user.Date); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, rows.Err()
}
