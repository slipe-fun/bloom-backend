package ServerMessageRepo

import "github.com/jmoiron/sqlx"

type ServerMessageRepo struct {
	db *sqlx.DB
}

func NewServerMessageRepo(db *sqlx.DB) *ServerMessageRepo {
	return &ServerMessageRepo{db: db}
}
