package chat

import (
	"github.com/jmoiron/sqlx"
	UserRepo "github.com/slipe-fun/skid-backend/internal/repository/user"
)

type ChatRepo struct {
	db       *sqlx.DB
	userRepo *UserRepo.UserRepo
}

func NewChatRepo(db *sqlx.DB, userRepo *UserRepo.UserRepo) *ChatRepo {
	return &ChatRepo{db: db, userRepo: userRepo}
}
