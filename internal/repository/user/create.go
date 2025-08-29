package UserRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *UserRepo) Create(user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (username, password) 
	          VALUES ($1, $2) 
	          RETURNING id, username, password, date`

	var created domain.User
	err := r.db.QueryRow(query, user.Username, user.Password).
		Scan(&created.ID, &created.Username, &created.Password, &created.Date)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
