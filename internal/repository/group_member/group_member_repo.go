package groupmember

import (
	"github.com/jmoiron/sqlx"
)

type GroupMemberRepo struct {
	db *sqlx.DB
}

func NewGroupMemberRepo(db *sqlx.DB) *GroupMemberRepo {
	return &GroupMemberRepo{db: db}
}
