package session

import (
	"fmt"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// LoginPacket ...
type LoginPacket struct{}

// Handle ...
func (p *LoginPacket) Handle(pk packet.Packet, conn *minecraft.Conn) error {
	login, ok := pk.(*packet.Login)
	if !ok {
		return fmt.Errorf("handle %T: expected: *packet.Login, got: %T", login, pk)
	}

	if login.ClientProtocol != protocol.CurrentProtocol {
		return fmt.Errorf("handle %T: unsupported protocol, expected: %v, got: %v", login, protocol.CurrentProtocol, login.ClientProtocol)
	}
	if len(login.ConnectionRequest) == 0 {
		return fmt.Errorf("handle %T: empty connection request", login)
	}

	return nil
}
