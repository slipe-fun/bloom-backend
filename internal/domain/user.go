package domain

import "time"

type User struct {
	ID       int       `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
	Date     time.Time `db:"date"`
}
