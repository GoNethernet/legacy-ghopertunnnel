package player

import "time"

// Title is a text that will appear in players screen.
type Title interface {
	// Text is the message that will be sent.
	Text() string
	// RemainDuration is the duration that the title will stay.
	RemainDuration() time.Duration
	// FadeInDuration is the fade in duration of the title.
	FadeInDuration() time.Duration
	// FadeOutDuration is the fade out duration of the title.
	FadeOutDuration() time.Duration
}

type title struct {
	text            string
	remainDuration  time.Duration
	fadeInDuration  time.Duration
	fadeOutDuration time.Duration
}

func (t *title) Text() string {
	return t.text
}
func (t *title) RemainDuration() time.Duration {
	return t.remainDuration
}
func (t *title) FadeInDuration() time.Duration {
	return t.fadeInDuration
}
func (t *title) FadeOutDuration() time.Duration {
	return t.fadeOutDuration
}
func newTitle(text string, remainDuration time.Duration, fadeInDuration time.Duration, fadeOutDuration time.Duration) *title {
	return &title{
		text:            text,
		remainDuration:  remainDuration,
		fadeInDuration:  fadeInDuration,
		fadeOutDuration: fadeOutDuration,
	}
}
