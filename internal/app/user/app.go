package user

type UserApp struct {
	users UserRepo
}

func NewUserApp(users UserRepo) *UserApp {
	return &UserApp{
		users: users,
	}
}
