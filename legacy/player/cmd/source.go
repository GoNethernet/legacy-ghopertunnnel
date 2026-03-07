package cmd

// Source is an interface that implements player.Player and defines the player that executes the command.
type Source interface {
	// RegisterCommand is a function contained in player.Player used for implementation.
	RegisterCommand(c Command)
}
