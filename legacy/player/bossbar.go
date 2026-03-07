package player

import "image/color"

// Bossbar a minecraft visual feature, in vanilla minecraft this is used for the wither, ender dragon and the raid.
type Bossbar interface {
	// Text is the text that will appear right up the health.
	Text() string
	// Health is the health of the bar.
	Health() int32
	// Colour represent that colour of the health bar.
	Colour() color.RGBA
}

type bossbar struct {
	text   string
	health int32
	colour color.RGBA
}

func (b *bossbar) Text() string {
	return b.text
}
func (b *bossbar) Health() int32 {
	return b.health
}
func (b *bossbar) Colour() color.RGBA {
	return b.colour
}
func newBossbar(text string, health int32, colour color.RGBA) *bossbar {
	return &bossbar{
		text:   text,
		health: health,
		colour: colour,
	}
}
