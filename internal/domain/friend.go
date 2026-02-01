package domain

type FriendRow struct {
	ID       int    `json:"id" db:"id"`
	UserID   int    `json:"user_id" db:"user_id"`
	FriendID int    `json:"friend_id" db:"friend_id"`
	Status   string `json:"status" db:"status"`
}

type Friend struct {
	ID       int    `json:"id" db:"id"`
	FriendID int    `json:"friend_id" db:"friend_id"`
	Status   string `json:"status" db:"status"`
}
