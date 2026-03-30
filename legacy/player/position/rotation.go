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

func RotationByValue(yaw, pitch float32) *Rotation {
	return &Rotation{
		pk: &packet.PlayerAuthInput{
			Yaw:   yaw,
			Pitch: pitch,
		},
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
	p.pk.Yaw = yaw
	return p.conn.WritePacket(&packet.MovePlayer{
		EntityRuntimeID: p.conn.GameData().EntityRuntimeID,
		Mode:            packet.MoveModeRotation,
		Yaw:             yaw,
		Pitch:           p.pk.Pitch,
	})
}

// SetPitch ...
func (p *Rotation) SetPitch(pitch float32) error {
	p.pk.Pitch = pitch
	return p.conn.WritePacket(&packet.MovePlayer{
		EntityRuntimeID: p.conn.GameData().EntityRuntimeID,
		Mode:            packet.MoveModeRotation,
		Pitch:           pitch,
		Yaw:             p.pk.Yaw,
	})
}

// String ...
func (p *Rotation) String() string {
	return fmt.Sprintf("yaw: %v, pitch: %v", p.Yaw(), p.Pitch())
}
