package session

import (
	"fmt"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// TextPacket ...
type TextPacket struct{}

// Handle ...
func (p *TextPacket) Handle(pk packet.Packet, conn *minecraft.Conn) error {
	text, ok := pk.(*packet.Text)
	if !ok {
		return fmt.Errorf("handle %T: expected: *packet.Text, got: %T", text, pk)
	}
	if text.SourceName != conn.IdentityData().DisplayName {
		if text.SourceName != "" {
			return fmt.Errorf("handle %T: mismatched names: %s and %s", text, text.SourceName, conn.IdentityData().DisplayName)
		}
	}
	if len(text.Message) >= 650 {
		return fmt.Errorf("handle %T: oversized text, len: %v", text, len(text.Message))
	}
	return nil
}
