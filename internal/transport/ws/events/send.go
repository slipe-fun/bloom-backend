package events

import "github.com/slipe-fun/skid-backend/internal/transport/ws/types"

func Send(hub *types.Hub, client *types.Client, room, message string) {
	hub.Broadcast(room, []byte(message))
}
