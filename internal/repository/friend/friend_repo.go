package FriendRepo

import "github.com/jmoiron/sqlx"

type FriendRepo struct {
	db *sqlx.DB
}

func NewFriendRepo(db *sqlx.DB) *FriendRepo {
	return &FriendRepo{db: db}
}
