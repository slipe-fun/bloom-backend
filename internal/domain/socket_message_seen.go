package domain

type SocketMessageSeen struct {
	Type     string `json:"type"`
	ChatID   int    `json:"chat_id"`
	Messages []int  `json:"messages"`
}
