package session

import (
	"sync"

	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/permission"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Session struct {
	mu             *sync.Mutex
	cpk, spk       packet.Packet
	health         int32
	gamemode       int32
	perm           int32
	heldSlot       int32
	inputs         *packet.PlayerAuthInput
	command        *packet.CommandRequest
	commandsPacket *packet.AvailableCommands
}

func New(mu *sync.Mutex) *Session {
	return &Session{mu: mu}
}

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
	}
}

func (s *Session) Health() int32   { return s.health }
func (s *Session) GameMode() int32 { return s.gamemode }

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

func (s *Session) HeldSlot() int32                 { return s.heldSlot }
func (s *Session) Inputs() *packet.PlayerAuthInput { return s.inputs }
func (s *Session) Command() *packet.CommandRequest { return s.command }
func (s *Session) ClearCommandData()               { s.command = nil }

func (s *Session) AvailableCommands() *packet.AvailableCommands {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.commandsPacket
}
