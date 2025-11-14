package domain

import "time"

type User struct {
	ID       int       `db:"id" json:"id"`
	Username string    `db:"username" json:"username"`
	Email    string    `db:"email" json:"email"`
	Date     time.Time `db:"date" json:"date"`
}
