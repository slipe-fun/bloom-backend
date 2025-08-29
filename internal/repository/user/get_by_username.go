package UserRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *UserRepo) GetByUsername(username string) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, username, password, date FROM users WHERE username = $1`
	err := r.db.Get(&user, query, username)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
