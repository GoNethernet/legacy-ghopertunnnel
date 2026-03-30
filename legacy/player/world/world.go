package world

import (
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/position"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/world/difficulty"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/world/dimension"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/world/gamemode"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/world/time"
	"github.com/sandertv/gophertunnel/minecraft"
)

// World serves as a high-level abstraction for environment-specific data,
// encapsulating state synchronization between the server's GameData and temporal packets.
type World struct {
	game minecraft.GameData
	time int32
}

// NewHandler initializes a new World controller utilizing the provided
// session metadata and initial temporal synchronization packet.
func NewHandler(p minecraft.GameData, time int32) *World {
	return &World{
		game: p,
		time: time,
	}
}

// Name returns the canonical string identifier of the world instance.
func (w *World) Name() string {
	return w.game.WorldName
}

// Difficulty evaluates the current environmental challenge tier,
// mapping protocol identifiers to concrete Difficulty implementations.
func (w *World) Difficulty() difficulty.Difficulty {
	if w.game.Hardcore {
		return difficulty.Hardcore{}
	}
	switch w.game.Difficulty {
	case 0:
		return difficulty.Peaceful{}
	case 1:
		return difficulty.Easy{}
	case 2:
		return difficulty.Normal{}
	case 3:
		return difficulty.Hard{}
	default:
		return nil
	}
}

// Seed retrieves the unique 64-bit numerical seed utilized for pseudo-random world generation.
func (w *World) Seed() int64 {
	return w.game.WorldSeed
}

// Dimension identifies the current spatial plane of existence within the world hierarchy.
func (w *World) Dimension() dimension.Dimension {
	switch w.game.Dimension {
	case 0:
		return dimension.Overworld{}
	case 1:
		return dimension.Nether{}
	case 2:
		return dimension.End{}
	default:
		return nil
	}
}

// Spawn returns the absolute block coordinates defined as the global entry point for the world.
func (w *World) Spawn() position.Position {
	var pos position.Position
	return pos.FromBlockPos(w.game.WorldSpawn)
}

// PlayerSpawn provides the specific vector coordinates designated for individual player manifestation.
func (w *World) PlayerSpawn() position.Position {
	var pos position.Position
	return pos.FromMgl32(w.game.PlayerPosition)
}

// GameMode returns the active interaction paradigm governing the player's capabilities and constraints.
func (w *World) GameMode() gamemode.GameMode {
	switch w.game.WorldGameMode {
	case 0:
		return gamemode.Survival{}
	case 1:
		return gamemode.Creative{}
	case 2:
		return gamemode.Adventure{}
	case 3:
		return gamemode.Spectator{}
	default:
		return nil
	}
}

// BaseGameVersion returns the specific protocol version string the world state is calibrated for.
func (w *World) BaseGameVersion() string {
	return w.game.BaseGameVersion
}

// Rotation extracts the player's initial angular orientation within the 3D spatial grid.
func (w *World) Rotation() *position.Rotation {
	return position.RotationByValue(w.game.Yaw, w.game.Pitch)
}

// Time evaluates the current diurnal state based on the synchronized tick value in the SetTime packet.
func (w *World) Time() time.Time {
	switch w.time {
	case 1000:
		return time.Day{}
	case 13000:
		return time.Night{}
	case 6000:
		return time.Noon{}
	case 18000:
		return time.Midnight{}
	case 23000:
		return time.Sunrise{}
	case 12000:
		return time.Sunset{}
	default:
		return time.Any{Ticks: int64(w.time)}
	}
}

// ExperimentsEnabled returns the operational status of experimental gameplay features within the world state.
func (w *World) ExperimentsEnabled() bool {
	return len(w.game.Experiments) > 0
}
