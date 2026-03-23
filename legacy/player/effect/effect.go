package effect

import "time"

// Effect ...
type Effect interface {
	// Type is the type of the effect, it is one of the constants above.
	Type() int32
	// Force is the effect amplifier.
	Force() float32
	// Duration defines the duration of the effect.
	Duration() time.Duration
	// Particles defines if the effect will show particles.
	Particles() bool
}

// All returns a slice of all effect names.
func All() []string {
	return []string{
		"speed", "slowness", "haste", "mining_fatigue", "strength",
		"instant_health", "instant_damage", "jump_boost", "nausea",
		"regeneration", "resistance", "fire_resistance", "water_breathing",
		"invisibility", "blindness", "night_vision", "hunger", "weakness",
		"poison", "wither", "health_boost", "absorption", "saturation",
		"levitation", "fatal_poison", "conduit_power", "slow_falling",
	}
}

// ByName returns an effect with its force and duration by its name.
func ByName(name string, force float32, duration time.Duration, particles bool) Effect {
	switch name {
	case "speed":
		return Speed(force, duration, particles)
	case "slowness":
		return Slowness(force, duration, particles)
	case "haste":
		return Haste(force, duration, particles)
	case "mining_fatigue":
		return MiningFatigue(force, duration, particles)
	case "strength":
		return Strength(force, duration, particles)
	case "instant_health":
		return InstantHealth(force, duration, particles)
	case "instant_damage":
		return InstantDamage(force, duration, particles)
	case "jump_boost":
		return JumpBoost(force, duration, particles)
	case "nausea":
		return Nausea(force, duration, particles)
	case "regeneration":
		return Regeneration(force, duration, particles)
	case "resistance":
		return Resistance(force, duration, particles)
	case "fire_resistance":
		return FireResistance(force, duration, particles)
	case "water_breathing":
		return WaterBreathing(force, duration, particles)
	case "invisibility":
		return Invisibility(force, duration, particles)
	case "blindness":
		return Blindness(force, duration, particles)
	case "night_vision":
		return NightVision(force, duration, particles)
	case "hunger":
		return Hunger(force, duration, particles)
	case "weakness":
		return Weakness(force, duration, particles)
	case "poison":
		return Poison(force, duration, particles)
	case "wither":
		return Wither(force, duration, particles)
	case "health_boost":
		return HealthBoost(force, duration, particles)
	case "absorption":
		return Absorption(force, duration, particles)
	case "saturation":
		return Saturation(force, duration, particles)
	case "levitation":
		return Levitation(force, duration, particles)
	case "fatal_poison":
		return FatalPoison(force, duration, particles)
	case "conduit_power":
		return ConduitPower(force, duration, particles)
	case "slow_falling":
		return SlowFalling(force, duration, particles)
	default:
		return nil
	}
}
