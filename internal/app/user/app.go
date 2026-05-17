package user

type UserApp struct {
	users UserRepo
	keys  KeysRepo
}

func NewUserApp(users UserRepo, keys KeysRepo) *UserApp {
	return &UserApp{
		users: users,
		keys:  keys,
	}
}
