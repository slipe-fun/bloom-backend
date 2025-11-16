package UserRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *UserRepo) GetByEmail(email string) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, username, email, display_name, date FROM users WHERE email = $1`
	err := r.db.Get(&user, query, email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
