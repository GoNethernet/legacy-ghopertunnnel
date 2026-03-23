package session

import (
	"sync"

	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/permission"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Session represent a single session in the proxy.
type Session struct {
	mu                               *sync.Mutex
	cpk, spk                         packet.Packet
	health, gamemode, perm, heldSlot int32
	hunger                           float32
	inputs                           *packet.PlayerAuthInput
	command                          *packet.CommandRequest
	commandsPacket                   *packet.AvailableCommands
}

// New creates a new empty session.
func New(mu *sync.Mutex) *Session {
	return &Session{mu: mu}
}

// UpdateFromClient updates packet from the client.
func (s *Session) UpdateFromClient(pk packet.Packet) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cpk = pk
	switch t := pk.(type) {
	case *packet.PlayerAuthInput:
		s.inputs = t
	case *packet.CommandRequest:
		s.command = t
	case *packet.MobEquipment:
		s.heldSlot = int32(t.InventorySlot)
	}
}

// UpdateFromServer updates packets from the server.
func (s *Session) UpdateFromServer(pk packet.Packet) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.spk = pk
	switch t := pk.(type) {
	case *packet.SetHealth:
		s.health = t.Health
	case *packet.SetPlayerGameType:
		s.gamemode = t.GameType
	case *packet.RequestPermissions:
		s.perm = t.PermissionLevel
	case *packet.AvailableCommands:
		s.commandsPacket = t
	case *packet.UpdateAttributes:
		for _, attr := range t.Attributes {
			if attr.Name == "minecraft:player.hunger" {
				s.hunger = attr.Value
			}
		}
	}
}

// Health ...
func (s *Session) Health() int32 { return s.health }

// GameMode ...
func (s *Session) GameMode() int32 { return s.gamemode }

// PermissionLevel ...
func (s *Session) PermissionLevel() permission.Permission {
	switch s.perm {
	case 0:
		return permission.Member{}
	case 1:
		return permission.Operator{}
	case 2:
		return permission.Visitor{}
	case 3:
		return permission.Custom{}
	default:
		return nil
	}
}

// HeldSlot ...
func (s *Session) HeldSlot() int32 { return s.heldSlot }

// Inputs returns player auth inputs.
func (s *Session) Inputs() *packet.PlayerAuthInput { return s.inputs }

// Command returns the commands request data packet.
func (s *Session) Command() *packet.CommandRequest { return s.command }

// ClearCommandData clears the command data.
func (s *Session) ClearCommandData() { s.command = nil }

// AvailableCommands ...
func (s *Session) AvailableCommands() *packet.AvailableCommands {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.commandsPacket
}

// Hunger ...
func (s *Session) Hunger() float32 {
	return s.hunger
}
