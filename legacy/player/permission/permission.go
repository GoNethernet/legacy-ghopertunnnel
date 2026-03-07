package permission

// LevelError returns an insufficient permission message.
func LevelError() string {
	return "§cYou do not have permission to use this command."
}

// Permission ...
type Permission interface {
	// Name ...
	Name() string
	// Level is an integer identifier for the permissions.
	Level() int32
}
