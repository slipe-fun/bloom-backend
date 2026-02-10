package auth

type AuthApp struct {
	sessionApp SessionApp
	users      UserRepo
	codesRepo  VerificationRepo
	codesApp   VerificationApp
	google     GoogleAuthService
}

func NewAuthApp(sessionApp SessionApp,
	users UserRepo,
	codesRepo VerificationRepo,
	codesApp VerificationApp,
	google GoogleAuthService) *AuthApp {
	return &AuthApp{
		sessionApp: sessionApp,
		users:      users,
		codesRepo:  codesRepo,
		codesApp:   codesApp,
		google:     google,
	}
}
