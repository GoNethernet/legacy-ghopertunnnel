package session

import (
	"fmt"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// RespawnPacket ...
type RespawnPacket struct{}

// Handle ...
func (r RespawnPacket) Handle(pk packet.Packet, conn *minecraft.Conn) error {
	respawn, ok := pk.(*packet.Respawn)
	if !ok {
		return fmt.Errorf("handle %T: expected: *packet.Respawn, got: %T", respawn, pk)
	}
	if &respawn.Position == nil {
		return fmt.Errorf("handle: %T: unknown respawn position", respawn)
	}
	return nil
}
