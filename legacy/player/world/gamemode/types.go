package gamemode

// Survival ...
type Survival struct{}

func (Survival) Name() string {
	return "survival"
}
func (Survival) ID() int32 {
	return 0
}
func (Survival) AllowsFlying() bool {
	return false
}

// Creative ...
type Creative struct{}

func (Creative) Name() string {
	return "creative"
}
func (Creative) ID() int32 {
	return 1
}
func (Creative) AllowsFlying() bool {
	return true
}

// Adventure ...
type Adventure struct{}

func (Adventure) Name() string {
	return "adventure"
}
func (Adventure) ID() int32 {
	return 2
}
func (Adventure) AllowsFlying() bool {
	return false
}

// Spectator ...
type Spectator struct{}

func (Spectator) Name() string {
	return "spectator"
}
func (Spectator) ID() int32 {
	return 3
}
func (Spectator) AllowsFlying() bool {
	return true
}
