package events

import "github.com/slipe-fun/skid-backend/internal/transport/ws/types"

func Join(hub *types.Hub, client *types.Client, room string) {
	hub.JoinRoom(client, room)
}
