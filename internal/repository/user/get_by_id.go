package UserRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *UserRepo) GetById(id int) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, username, password, date FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
