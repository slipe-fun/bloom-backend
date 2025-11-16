package UserRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *UserRepo) Edit(user *domain.User) error {
	query := `
        UPDATE users
        SET username = $1,
            email = $2,
			display_name = $3
        WHERE id = $4
        RETURNING id, username, email, display_name, date
    `

	_, err := r.db.Exec(query, user.Username, user.Email, user.DisplayName, user.ID)

	if err != nil {
		return err
	}

	return nil
}
