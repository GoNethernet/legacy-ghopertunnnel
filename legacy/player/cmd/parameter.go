package cmd

// Optional represents a parameter that may or may not be provided in a command.
type Optional[T any] struct {
	// value holds the actual value of the parameter.
	value T
	// Has indicates if the parameter was actually provided by the sender.
	Has bool
}

// Value ...
func (o Optional[T]) Value() T {
	return o.value
}

// Enum represents a command parameter with a fixed set of options.
type Enum interface {
	// Type returns the name of the enum type.
	Type() string
	// Options returns a list of all valid values for the enum.
	Options() []string
}

// SubCommand represents a literal argument for sub-commands.
type SubCommand struct{}

// Varargs represents a string that consumes all remaining arguments.
type Varargs string

// Target represents a single entity or player targeted by a command.
type Target interface {
	// Name returns the name of the target.
	Name() string
}

// PlayerTarget is a basic implementation of Target for players.
type PlayerTarget struct {
	NameValue string
}

func (t PlayerTarget) Name() string { return t.NameValue }
