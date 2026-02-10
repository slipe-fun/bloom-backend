package server

import "github.com/jmoiron/sqlx"

type ServerChannelRepo struct {
	db *sqlx.DB
}

func NewServerChannelRepo(db *sqlx.DB) *ServerChannelRepo {
	return &ServerChannelRepo{db: db}
}
