package effect

import "time"

// Effect ...
type Effect interface {
	// Type is the type of the effect, it is one of the constants above.
	Type() int32
	// Force is the effect amplifier.
	Force() float32
	// Duration defines the duration of the effect.
	Duration() time.Duration
	// Particles defines if the effect will show particles.
	Particles() bool
}
