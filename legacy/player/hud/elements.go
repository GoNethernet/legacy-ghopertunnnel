package hud

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

// PaperDoll is the element that shows the player's paper doll, which is a visual representation of the
// player's character model and equipment, as well as any currently played animations. It is located in the
// top left corner of the screen.
type PaperDoll struct{}

func (p PaperDoll) Name() string { return "PaperDoll" }
func (p PaperDoll) Type() int32  { return packet.HudElementPaperDoll }

// Armour is the element that shows the player's armour level, sitting either above the hotbar or at the top
// of the screen on in non-classic views.
type Armour struct{}

func (a Armour) Name() string { return "Armour" }
func (a Armour) Type() int32  { return packet.HudElementArmour }

// ToolTips is the element that shows useful hints and tips to the player, such as how to use items or
// how to perform certain actions in the game. These tips are displayed at the top right of the screen.
type ToolTips struct{}

func (t ToolTips) Name() string { return "ToolTips" }
func (t ToolTips) Type() int32  { return packet.HudElementToolTips }

// TouchControls is the element that shows the touch controls on the screen, which is used for touch-based
// devices.
type TouchControls struct{}

func (t TouchControls) Name() string { return "TouchControls" }
func (t TouchControls) Type() int32  { return packet.HudElementTouchControls }

// Crosshair is the element that shows the crosshair in the middle of the screen, which is used for aiming
// and targeting entities or blocks.
type Crosshair struct{}

func (c Crosshair) Name() string { return "Crosshair" }
func (c Crosshair) Type() int32  { return packet.HudElementCrosshair }

// HotBar is the element that shows all the items in the player's hotbar, located at the bottom of the screen.
type HotBar struct{}

func (h HotBar) Name() string { return "HotBar" }
func (h HotBar) Type() int32  { return packet.HudElementHotBar }

// Health is the element that shows the player's health bar, sitting either above the hotbar or at the top
// of the screen on in non-classic views.
type Health struct{}

func (h Health) Name() string { return "Health" }
func (h Health) Type() int32  { return packet.HudElementHealth }

// ProgressBar is the element that shows the player's experience bar. It is always located just above the
// hotbar.
type ProgressBar struct{}

func (p ProgressBar) Name() string { return "ProgressBar" }
func (p ProgressBar) Type() int32  { return packet.HudElementProgressBar }

// Hunger is the element that shows the player's hunger bar, which indicates how hungry the player is and
// how much food they need to consume to restore their hunger. It is located either above the hotbar or at the
// top of the screen on in non-classic views.
type Hunger struct{}

func (h Hunger) Name() string { return "Hunger" }
func (h Hunger) Type() int32  { return packet.HudElementHunger }

// AirBubbles is the element that shows the player's air bubbles, which indicate how much air the player has
// left when underwater. It is located either above the hotbar or at the top of the screen on in non-classic
// views. It is only visible when the player is underwater or they are regenerating air after being underwater.
type AirBubbles struct{}

func (a AirBubbles) Name() string { return "AirBubbles" }
func (a AirBubbles) Type() int32  { return packet.HudElementAirBubbles }

// HorseHealth is the element that shows the health of the player's horse, which replaces the player's own
// health bar when riding a horse/other entity with health.
type HorseHealth struct{}

func (h HorseHealth) Name() string { return "HorseHealth" }
func (h HorseHealth) Type() int32  { return packet.HudElementHorseHealth }

// StatusEffects is the element that shows the icons of the currently active status effects, located on the
// right side of the screen.
type StatusEffects struct{}

func (s StatusEffects) Name() string { return "StatusEffects" }
func (s StatusEffects) Type() int32  { return packet.HudElementStatusEffects }

// ItemText is the element that shows the text of the item currently held in the player's hand, which is
// displayed just above the hotbar when switching to a new item.
type ItemText struct{}

func (i ItemText) Name() string { return "ItemText" }
func (i ItemText) Type() int32  { return packet.HudElementItemText }

// All returns all the HUD elements that are available to be shown or hidden in the game.
func All() []Hud {
	return []Hud{
		PaperDoll{}, Armour{}, ToolTips{}, TouchControls{}, Crosshair{}, HotBar{}, Health{},
		ProgressBar{}, Hunger{}, AirBubbles{}, HorseHealth{}, StatusEffects{}, ItemText{},
	}
}
