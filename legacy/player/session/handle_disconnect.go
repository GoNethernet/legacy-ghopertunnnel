package session

import (
	"fmt"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// DisconnectPacket ...
type DisconnectPacket struct{}

// Handle ...
func (p *DisconnectPacket) Handle(pk packet.Packet, conn *minecraft.Conn) error {
	dis, ok := pk.(*packet.Disconnect)
	if !ok {
		return fmt.Errorf("handle %T: expected: *packet.Disconnect, got: %T", dis, pk)
	}
	if dis.Message == "" {
		return fmt.Errorf("handle %T: disconnected with no message provided", dis)
	}
	if dis.Reason == 0 {
		return fmt.Errorf("handle %T: disconnected with no reason provided", dis)
	}
	return nil
}
