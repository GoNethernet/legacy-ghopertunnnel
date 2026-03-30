package gamemode

// GameMode defines the behavioral interface for player interaction states,
// governing environmental engagement, protocol identification, and aerial mobility permissions.
type GameMode interface {
	// Name returns the canonical nomenclature of the interaction paradigm.
	Name() string
	// ID returns the 32-bit signed integer identifier utilized for protocol synchronization.
	ID() int32
	// AllowsFlying indicates whether the interaction state permits self-propelled aerial navigation.
	AllowsFlying() bool
}
