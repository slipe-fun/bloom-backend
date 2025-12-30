package UserRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *UserRepo) Create(user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (username, email, display_name, description) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id, username, email, display_name, description, date`

	var created domain.User
	err := r.db.QueryRow(query, user.Username, user.Email, user.DisplayName, user.Description).
		Scan(&created.ID, &created.Username, &created.Email, &created.DisplayName, &created.Description, &created.Date)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
