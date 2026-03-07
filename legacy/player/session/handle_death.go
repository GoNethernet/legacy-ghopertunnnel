package session

import (
	"fmt"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// DeathPacket ...
type DeathPacket struct{}

// Handle ...
func (p DeathPacket) Handle(pk packet.Packet, conn *minecraft.Conn) error {
	death, ok := pk.(*packet.DeathInfo)
	if !ok {
		return fmt.Errorf("handle %T: expected: *packet.DeathInfo, got: %T", p, pk)
	}
	if death.Cause == "" {
		return fmt.Errorf("handle %T: no cause found", death)
	}
	return nil
}
