package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

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

func (r *UserRepo) GetByUsername(username string) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, username, password, date FROM users WHERE username = $1`
	err := r.db.Get(&user, query, username)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetById(id int) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, username, password, date FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
