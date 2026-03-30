package difficulty

// Peaceful ...
type Peaceful struct{}

func (Peaceful) Name() string {
	return "peaceful"
}
func (Peaceful) ID() uint64 {
	return 0
}

// Easy ...
type Easy struct{}

func (Easy) Name() string {
	return "easy"
}
func (Easy) ID() uint64 {
	return 1
}

// Normal ...
type Normal struct{}

func (Normal) Name() string {
	return "normal"
}
func (Normal) ID() uint64 {
	return 2
}

// Hard ...
type Hard struct{}

func (Hard) Name() string {
	return "hard"
}
func (Hard) ID() uint64 {
	return 3
}

// Hardcore ...
type Hardcore struct{}

func (Hardcore) Name() string {
	return "hardcore"
}
func (Hardcore) ID() uint64 {
	return 4
}
