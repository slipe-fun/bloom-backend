package credential

import (
	"github.com/jmoiron/sqlx"
)

type CredentialRepo struct {
	db *sqlx.DB
}

func NewCredentialRepo(db *sqlx.DB) *CredentialRepo {
	return &CredentialRepo{db: db}
}
