package cmd

// optionalT is an additional interface implementation for Optional.
type optionalT interface {
	// with represent an Optional type that addresses a value.
	with(val any) any
}

// Optional represents a parameter that may or may not be provided in a command.
type Optional[T any] struct {
	// value holds the actual value of the parameter.
	Value T
	// set indicates if the parameter was actually provided by the sender.
	Has bool
}

// Load returns the underlying value and a boolean indicating if the parameter was provided.
func (o Optional[T]) Load() (T, bool) {
	return o.Value, o.Has
}

// LoadOr returns the value if the parameter was provided, or a default value 'or' if it was not.
func (o Optional[T]) LoadOr(or T) T {
	if o.Has {
		return o.Value
	}
	return or
}
func (o Optional[T]) with(val any) any {
	return Optional[T]{Value: val.(T), Has: true}
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
