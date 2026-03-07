package session

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Handler ...
type Handler interface {
	// Handle ...
	Handle(pk packet.Packet, conn *minecraft.Conn) error
}
