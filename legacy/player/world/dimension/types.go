package dimension

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

// Overworld ...
type Overworld struct{}

func (Overworld) Name() string {
	return "overworld"
}
func (Overworld) ID() uint32 {
	return packet.DimensionOverworld
}

// Nether ...
type Nether struct{}

func (Nether) Name() string {
	return "nether"
}
func (Nether) ID() uint32 {
	return packet.DimensionNether
}

// End ...
type End struct{}

func (End) Name() string {
	return "end"
}
func (End) ID() uint32 {
	return packet.DimensionEnd
}
