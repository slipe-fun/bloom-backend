package server

import "github.com/jmoiron/sqlx"

type ServerRepo struct {
	db *sqlx.DB
}

func NewServerRepo(db *sqlx.DB) *ServerRepo {
	return &ServerRepo{db: db}
}
