package user

type UserHandler struct {
	userApp   UserApp
}

func NewUserHandler(userApp UserApp) *UserHandler {
	return &UserHandler{
		userApp:   userApp,
	}
}
