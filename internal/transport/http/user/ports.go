package user

import "github.com/slipe-fun/skid-backend/internal/domain"

type UserApp interface {
	GetUserByToken(token string) (*domain.User, error)
	GetUserByID(id int) (*domain.User, error)
	SearchUsersByUsername(username string, limit, offset int) ([]*domain.User, error)
	IsUserWithEmailExists(email string) (bool, error)
	EditUser(token string, editedUser *domain.User) (*domain.User, error)
}

type FriendApp interface {
	GetFriendCount(token string) (int, error)
}
