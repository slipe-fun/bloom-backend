package user

import (
	"github.com/jmoiron/sqlx"
	VerificationRepo "github.com/slipe-fun/skid-backend/internal/repository/verification"
)

type UserRepo struct {
	db               *sqlx.DB
	verificationRepo *VerificationRepo.VerificationRepo
}

func NewUserRepo(db *sqlx.DB, verificationRepo *VerificationRepo.VerificationRepo) *UserRepo {
	return &UserRepo{db: db, verificationRepo: verificationRepo}
}
