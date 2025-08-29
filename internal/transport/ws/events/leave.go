package events

import "github.com/slipe-fun/skid-backend/internal/transport/ws/types"

func Leave(hub *types.Hub, client *types.Client) {
	hub.LeaveRoom(client)
}
