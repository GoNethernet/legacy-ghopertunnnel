package hud

// Hud represents a HUD element in the game that can either be hidden or shown.
type Hud interface {
	// Name ...
	Name() string
	// Type ...
	Type() int32
}
