package repository

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

type ChatRepo struct {
	db *sqlx.DB
}

func NewChatRepo(db *sqlx.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

func (r *ChatRepo) Create(chat *domain.Chat) error {
	membersJSON, _ := json.Marshal(chat.Members)
	query := `INSERT INTO chats (members) VALUES ($1) RETURNING id`
	return r.db.QueryRow(query, membersJSON).Scan(&chat.ID)
}

func (r *ChatRepo) GetById(id int) (*domain.Chat, error) {
	var chat domain.Chat
	var membersJSON []byte

	query := `SELECT id, members FROM chats WHERE id=$1`
	err := r.db.QueryRow(query, id).Scan(&chat.ID, &membersJSON)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(membersJSON, &chat.Members)

	return &chat, nil
}
