package server

import "github.com/jmoiron/sqlx"

type ServerMemberRepo struct {
	db *sqlx.DB
}

func NewServerMemberRepo(db *sqlx.DB) *ServerMemberRepo {
	return &ServerMemberRepo{db: db}
}
