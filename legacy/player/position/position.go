package position

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// Position ...
type Position [3]float64

func NewPosition(x, y, z float64) Position {
	return Position{x, y, z}
}

// X ...
func (p *Position) X() float64 {
	return p[0]
}

// Y ...
func (p *Position) Y() float64 {
	return p[1]
}

// Z ...
func (p *Position) Z() float64 {
	return p[2]
}

// ToMgl32 ...
func (p *Position) ToMgl32() mgl32.Vec3 {
	return mgl32.Vec3{float32(p.X()), float32(p.Y()), float32(p.Z())}
}

// ToMgl64 ...
func (p *Position) ToMgl64() mgl64.Vec3 {
	return mgl64.Vec3{p.X(), p.Y(), p.Z()}
}

// ToBlockPos ...
func (p *Position) ToBlockPos() protocol.BlockPos {
	return protocol.BlockPos{int32(p.X()), int32(p.Y()), int32(p.Z())}
}

// ToMgl32Pointer ...
func (p *Position) ToMgl32Pointer() *mgl32.Vec3 {
	return &mgl32.Vec3{float32(p.X()), float32(p.Y()), float32(p.Z())}
}

// ToMgl64Pointer ...
func (p *Position) ToMgl64Pointer() *mgl64.Vec3 {
	return &mgl64.Vec3{p.X(), p.Y(), p.Z()}
}

// ToBlockPosPointer ...
func (p *Position) ToBlockPosPointer() *protocol.BlockPos {
	return &protocol.BlockPos{int32(p.X()), int32(p.Y()), int32(p.Z())}
}

// FromBlockPos ...
func (p *Position) FromBlockPos(pos protocol.BlockPos) Position {
	return Position{float64(pos.X()), float64(pos.Y()), float64(pos.Z())}
}

// FromMgl32 ...
func (p *Position) FromMgl32(pos mgl32.Vec3) Position {
	return Position{float64(pos.X()), float64(pos.Y()), float64(pos.Z())}
}

// FromMgl64 ...
func (p *Position) FromMgl64(pos mgl64.Vec3) Position {
	return Position{pos.X(), pos.Y(), pos.Z()}
}

// FromMgl32Pointer ...
func (p *Position) FromMgl32Pointer(pos mgl32.Vec3) *Position {
	return &Position{float64(pos.X()), float64(pos.Y()), float64(pos.Z())}
}

// FromMgl64Pointer ...
func (p *Position) FromMgl64Pointer(pos *mgl64.Vec3) *Position {
	return &Position{pos.X(), pos.Y(), pos.Z()}
}

// FromBlockPosPointer ...
func (p *Position) FromBlockPosPointer(pos protocol.BlockPos) *Position {
	return &Position{float64(pos.X()), float64(pos.Y()), float64(pos.Z())}
}

// String ...
func (p *Position) String() string {
	return fmt.Sprintf("%f ,%f ,%f", p.X(), p.Y(), p.Z())
}
