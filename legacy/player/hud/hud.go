package hud

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

// Hud represents a HUD element in the game that can either be hidden or shown.
type Hud interface {
	// Name ...
	Name() string
	// Type ...
	Type() int32
}

// hudElement ...
type hudElement struct {
	name string
	id   int32
}

// Name ...
func (h hudElement) Name() string { return h.name }

// Type ...
func (h hudElement) Type() int32 { return h.id }

// All returns all the HUD elements that are available to be shown or hidden in the game.
func All() []Hud {
	return []Hud{
		hudElement{"paper_doll", packet.HudElementPaperDoll},
		hudElement{"armour", packet.HudElementArmour},
		hudElement{"tool_tips", packet.HudElementToolTips},
		hudElement{"touch_controls", packet.HudElementTouchControls},
		hudElement{"crosshair", packet.HudElementCrosshair},
		hudElement{"hotbar", packet.HudElementHotBar},
		hudElement{"health", packet.HudElementHealth},
		hudElement{"progress_bar", packet.HudElementProgressBar},
		hudElement{"hunger", packet.HudElementHunger},
		hudElement{"air_bubbles", packet.HudElementAirBubbles},
		hudElement{"horse_health", packet.HudElementHorseHealth},
		hudElement{"status_effects", packet.HudElementStatusEffects},
		hudElement{"item_text", packet.HudElementItemText},
	}
}

// Names ...
func Names() []string {
	all := All()
	names := make([]string, 0, len(all))
	for _, h := range all {
		names = append(names, h.Name())
	}
	return names
}

// ByName lookup for a hud with a name.
func ByName(name string) Hud {
	for _, h := range All() {
		if h.Name() == name {
			return h
		}
	}
	return nil
}
