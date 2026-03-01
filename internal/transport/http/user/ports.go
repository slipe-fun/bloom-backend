package user

import "github.com/slipe-fun/skid-backend/internal/domain"

type UserApp interface {
	GetAllUsers(limit, offset int) ([]*domain.User, error)
	GetUserByID(id int) (*domain.User, error)
	SearchUsersByUsername(username string, limit, offset int) ([]*domain.User, error)
	DoesUserExistGivenEmail(email string) (bool, error)
	EditUser(user_id int, editedUser *domain.User) (*domain.User, error)
}

type FriendApp interface {
	GetFriendCount(user_id int) (int, error)
}
