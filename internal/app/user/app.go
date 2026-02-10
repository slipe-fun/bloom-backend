package user

type UserApp struct {
	sessionApp SessionApp
	users      UserRepo
}

func NewUserApp(sessionApp SessionApp,
	users UserRepo,
) *UserApp {
	return &UserApp{
		sessionApp: sessionApp,
		users:      users,
	}
}
