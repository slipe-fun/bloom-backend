package auth

type AuthApp struct {
	sessionApp SessionApp
	users      UserRepo
}

func NewAuthApp(
	sessionApp SessionApp,
	users UserRepo,
) *AuthApp {
	return &AuthApp{
		sessionApp: sessionApp,
		users:      users,
	}
}
