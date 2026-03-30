package time

// Time defines the fundamental interface for temporal state representation tick-based.
type Time interface {
	// Name ...
	Name() string
	// ID returns the 64-bit integer representing the specific tick value
	// within the diurnal cycle.
	ID() int64
}
