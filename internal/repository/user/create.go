package UserRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *UserRepo) Create(user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (username, email, display_name) 
	          VALUES ($1, $2, $3) 
	          RETURNING id, username, email, display_name, date`

	var created domain.User
	err := r.db.QueryRow(query, user.Username, user.Email, user.DisplayName).
		Scan(&created.ID, &created.Username, &created.Email, &created.DisplayName, &created.Date)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
