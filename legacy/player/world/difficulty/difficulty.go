package difficulty

// Difficulty defines the foundational interface for environmental challenge tiers,
// encapsulating the operational parameters and numerical identification within the world state.
type Difficulty interface {
	// Name returns the canonical string representation of the difficulty level.
	Name() string
	// ID returns the unique 64-bit integer identifier, constrained to a range between 0 and 4.
	ID() uint64
}
