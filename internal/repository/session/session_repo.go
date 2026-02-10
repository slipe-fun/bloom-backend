package session

import (
	"github.com/jmoiron/sqlx"
	UserRepo "github.com/slipe-fun/skid-backend/internal/repository/user"
)

type SessionRepo struct {
	db *sqlx.DB
}

func NewSessionRepo(db *sqlx.DB, userRepo *UserRepo.UserRepo) *SessionRepo {
	return &SessionRepo{db: db}
}
