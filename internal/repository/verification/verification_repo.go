package verification

import "github.com/jmoiron/sqlx"

type VerificationRepo struct {
	db *sqlx.DB
}

func NewVerificationRepo(db *sqlx.DB) *VerificationRepo {
	return &VerificationRepo{db: db}
}
