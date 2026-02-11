package chat

import "github.com/slipe-fun/skid-backend/internal/domain"

type UserRepo interface {
	GetByID(id int) (*domain.User, error)
}
