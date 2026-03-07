package position

import (
	"fmt"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Rotation ...
type Rotation struct {
	conn *minecraft.Conn
	pk   *packet.PlayerAuthInput
}

// NewRotation ...
func NewRotation(conn *minecraft.Conn, pk *packet.PlayerAuthInput) *Rotation {
	return &Rotation{
		conn: conn,
		pk:   pk,
	}
}

// Yaw is the vertical rotation of the player.
func (p *Rotation) Yaw() float32 {
	return p.pk.Yaw
}

// Pitch is the horizontal rotation of the player.
func (p *Rotation) Pitch() float32 {
	return p.pk.Pitch
}

// SetYaw ...
func (p *Rotation) SetYaw(yaw float32) error {
	return p.conn.WritePacket(&packet.MovePlayer{
		EntityRuntimeID: p.conn.GameData().EntityRuntimeID,
		Mode:            packet.MoveModeRotation,
		Yaw:             yaw,
	})
}

// SetPitch ...
func (p *Rotation) SetPitch(pitch float32) error {
	return p.conn.WritePacket(&packet.MovePlayer{
		EntityRuntimeID: p.conn.GameData().EntityRuntimeID,
		Mode:            packet.MoveModeRotation,
		Pitch:           pitch,
	})
}

// String ...
func (p *Rotation) String() string {
	return fmt.Sprintf("yaw: %v, pitch: %v", p.Yaw(), p.Pitch())
}
