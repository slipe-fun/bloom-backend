package friend

import (
	"github.com/slipe-fun/skid-backend/internal/app/session"
	"github.com/slipe-fun/skid-backend/internal/auth"
	"github.com/slipe-fun/skid-backend/internal/repository/friend"
	"github.com/slipe-fun/skid-backend/internal/repository/user"
)

type FriendApp struct {
	sessionApp *session.SessionApp
	friends    *friend.FriendRepo
	users      *user.UserRepo
	tokenSvc   *auth.TokenService
}

func NewFriendApp(sessionApp *session.SessionApp,
	friends *friend.FriendRepo,
	users *user.UserRepo,
	tokenSvc *auth.TokenService) *FriendApp {
	return &FriendApp{
		sessionApp: sessionApp,
		friends:    friends,
		users:      users,
		tokenSvc:   tokenSvc,
	}
}
