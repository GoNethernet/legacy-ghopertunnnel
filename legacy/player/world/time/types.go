package time

// Day ...
type Day struct{}

func (Day) Name() string {
	return "day"
}
func (Day) ID() int64 {
	return 1000
}

// Night ...
type Night struct{}

func (Night) Name() string {
	return "night"
}
func (Night) ID() int64 {
	return 13000
}

// Noon ...
type Noon struct{}

func (Noon) Name() string {
	return "noon"
}
func (Noon) ID() int64 {
	return 6000
}

// Midnight ...
type Midnight struct{}

func (Midnight) Name() string {
	return "midnight"
}
func (Midnight) ID() int64 {
	return 18000
}

// Sunrise ...
type Sunrise struct{}

func (Sunrise) Name() string {
	return "sunrise"
}
func (Sunrise) ID() int64 {
	return 23000
}

// Sunset ...
type Sunset struct{}

func (Sunset) Name() string {
	return "sunset"
}
func (Sunset) ID() int64 {
	return 12000
}

// Any ...
type Any struct{ Ticks int64 }

func (Any) Name() string {
	return "unknown"
}
func (p Any) ID() int64 {
	return p.Ticks
}
