package domain

import "time"

type User struct {
	ID          int       `db:"id" json:"id"`
	Username    string    `db:"username" json:"username"`
	Email       *string   `db:"email" json:"email"`
	DisplayName *string   `db:"display_name" json:"display_name"`
	Description *string   `db:"description" json:"description"`
	Date        time.Time `db:"date" json:"date"`
}
