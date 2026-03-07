package effect

import (
	"time"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// effect ...
type effect struct {
	id        int32
	force     float32
	duration  time.Duration
	particles bool
}

func (e effect) Type() int32             { return e.id }
func (e effect) Force() float32          { return e.force }
func (e effect) Duration() time.Duration { return e.duration }
func (e effect) Particles() bool         { return e.particles }

// New returns a new effect with its properties.
func New(id int32, force float32, duration time.Duration, particles bool) Effect {
	return effect{id: id, force: force, duration: duration, particles: particles}
}

// Blindness ...
func Blindness(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectBlindness, force, duration, particles)
}

// Speed ...
func Speed(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectSpeed, force, duration, particles)
}

// Weakness ...
func Weakness(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectWeakness, force, duration, particles)
}

// Wither ...
func Wither(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectWither, force, duration, particles)
}

// Poison ...
func Poison(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectPoison, force, duration, particles)
}

// JumpBoost ...
func JumpBoost(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectJumpBoost, force, duration, particles)
}

// MiningFatigue ...
func MiningFatigue(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectMiningFatigue, force, duration, particles)
}

// FatalPoison ...
func FatalPoison(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectFatalPoison, force, duration, particles)
}

// Strength ...
func Strength(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectStrength, force, duration, particles)
}

// WaterBreathing ...
func WaterBreathing(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectWaterBreathing, force, duration, particles)
}

// Haste ...
func Haste(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectHaste, force, duration, particles)
}

// Resistance ...
func Resistance(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectResistance, force, duration, particles)
}

// Invisibility ...
func Invisibility(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectInvisibility, force, duration, particles)
}

// NightVision ...
func NightVision(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectNightVision, force, duration, particles)
}

// Absorption ...
func Absorption(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectAbsorption, force, duration, particles)
}

// ConduitPower ...
func ConduitPower(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectConduitPower, force, duration, particles)
}

// FireResistance ...
func FireResistance(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectFireResistance, force, duration, particles)
}

// HealthBoost ...
func HealthBoost(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectHealthBoost, force, duration, particles)
}

// Hunger ...
func Hunger(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectHunger, force, duration, particles)
}

// InstantDamage ...
func InstantDamage(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectInstantDamage, force, duration, particles)
}

// InstantHealth ...
func InstantHealth(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectInstantHealth, force, duration, particles)
}

// Levitation ...
func Levitation(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectLevitation, force, duration, particles)
}

// Nausea ...
func Nausea(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectNausea, force, duration, particles)
}

// Regeneration ...
func Regeneration(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectRegeneration, force, duration, particles)
}

// Saturation ...
func Saturation(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectSaturation, force, duration, particles)
}

// SlowFalling ...
func SlowFalling(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectSlowFalling, force, duration, particles)
}

// Slowness ...
func Slowness(force float32, duration time.Duration, particles bool) Effect {
	return New(packet.EffectSlowness, force, duration, particles)
}
