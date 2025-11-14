package UserRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *UserRepo) Create(user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (username, email) 
	          VALUES ($1, $2) 
	          RETURNING id, username, email, date`

	var created domain.User
	err := r.db.QueryRow(query, user.Username, user.Email).
		Scan(&created.ID, &created.Username, &created.Email, &created.Date)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
