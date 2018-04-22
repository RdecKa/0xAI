package hex

import "fmt"

// ------------------
// |     Action     |
// ------------------

// Action shows where next stone will be placed and of which color
type Action struct {
	x, y byte
	c    color
}

// NewAction creates a new action. x and y are coordinates of a stone placed by
// the player c
func NewAction(x, y byte, c color) *Action {
	return &Action{x, y, c}
}

func (a Action) String() string {
	return fmt.Sprintf("%s: (%d, %d)", a.c, a.x, a.y)
}
