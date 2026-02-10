package friend

type FriendApp struct {
	sessionApp SessionApp
	friends    FriendRepo
	users      UserRepo
}

func NewFriendApp(sessionApp SessionApp,
	friends FriendRepo,
	users UserRepo,
) *FriendApp {
	return &FriendApp{
		sessionApp: sessionApp,
		friends:    friends,
		users:      users,
	}
}
