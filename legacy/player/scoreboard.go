package player

// Scoreboard ...
type Scoreboard interface {
	// Title is the title of the scoreboard.
	Title() string
	// EntryID is the entry id of the scoreboard, the number is not important but has to be
	// different for each scoreboard.
	EntryID() int64
	// Lines are all the lines in the scoreboard, it doesn't include the title.
	Lines() []string
}
type scoreboard struct {
	title string
	id    int64
	lines []string
}

func (p *scoreboard) Title() string {
	return p.title
}
func (p *scoreboard) EntryID() int64 {
	return p.id
}
func (p *scoreboard) Lines() []string {
	return p.lines
}
func newScoreboard(title string, id int64, lines []string) *scoreboard {
	return &scoreboard{
		title: title,
		id:    id,
		lines: lines,
	}
}
