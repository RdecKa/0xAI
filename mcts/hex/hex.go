package hex

import (
	"fmt"
)

// -----------------
// |     Color     |
// -----------------

type color byte

// enum for players (colors)
const (
	None color = 0
	Red  color = 1
	Blue color = 2
)

func (c color) String() string {
	switch c {
	case Red:
		return "r"
	case Blue:
		return "b"
	default:
		return "."
	}
}

func getColorFromBits(bits uint64) color {
	if bits == 1 {
		return Red
	} else if bits == 2 {
		return Blue
	}
	return None
}

// -----------------
// |     State     |
// -----------------

// State represents a state on a grid in a hex game
//	size is a length of the grid (size 11 means 11x11 grid)
//	grid is a list of rows in a grid, each row is represented as uint64. Each
//		cell in a row is stored with two bits:
//			00 - empty
//			01 - red
//			10 - blue
//			11 - undefined
//		Lowest two bits represent the cell with index 0. Because of using 64
//		bits for a row, maximal size of the grid is 32x32.
//	lastPlayer is the color of the player who made the last action
type State struct {
	size       byte
	grid       []uint64
	lastPlayer color
}

// NewState returns new State with a grid of given size
func NewState(size byte) *State {
	grid := make([]uint64, size)
	return &State{size, grid, Blue} // Blue is set as last player, so Red always starts
}

func (s State) String() string {
	r := ""
	for rowIndex, row := range s.grid {
		for i := 0; i < rowIndex; i++ {
			r += " "
		}
		for col := byte(0); col < s.size; col++ {
			color := getCellInRow(row, col)
			r += fmt.Sprintf("%s ", color)
		}
		r += "\n"
	}
	return r
}

// getColorOn returns the color of the stone in cell (x, y)
func (s State) getColorOn(x, y byte) color {
	row := s.grid[y]
	return getCellInRow(row, x)
}

// setCell puts a stone of color c into cell (x, y)
// Cell (x, y) must be empty and valid
func (s *State) setCell(x, y byte, c color) {
	bits := uint64(c << (x * 2))
	s.grid[y] |= bits
}

// getCellInRow returns color of a stone on index index in row row
func getCellInRow(row uint64, index byte) color {
	// Find the two bits that represent column with index index
	bits := ((3 << (index * 2)) & row) >> (index * 2)
	return getColorFromBits(bits)
}

func (s State) clone() *State {
	return &State{s.size, s.grid, s.lastPlayer}
}

// GetSuccessorState returns a state after Action a is performed
func (s State) GetSuccessorState(a Action) *State {
	newState := s.clone()
	if a.c == s.lastPlayer {
		panic(fmt.Sprintf("Player cannot do two moves in a row! (last player: %s, current action: %s)", s.lastPlayer, a.c))
	}
	newState.lastPlayer = a.c
	newState.setCell(a.x, a.y, a.c)
	return newState
}

// ------------------
// |     Action     |
// ------------------

// Action shows where next stone will be placed and of which color
type Action struct {
	x, y byte
	c    color
}

// NewAction creates a new action. x and y are coordinates of a stone placed by
// color color
func NewAction(x, y byte, c color) *Action {
	return &Action{x, y, c}
}

func (a Action) String() string {
	return fmt.Sprintf("%s: (%d, %d)", a.c, a.x, a.y)
}
