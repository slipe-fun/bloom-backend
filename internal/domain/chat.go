package domain

type Chat struct {
	ID      int      `db:"id"`
	Members []string `db:"members"`
}
