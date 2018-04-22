package hex

import "fmt"

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
//
// A goal of the red player is to connect top-most and bottom-most row while a
// goal of the blue player is to connect left-most and right-most column
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

func (s *State) String() string {
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

// GetSize returns size of the grid
func (s *State) GetSize() int {
	return int(s.size)
}

// getColorOn returns the color of the stone in cell (x, y)
func (s *State) getColorOn(x, y byte) color {
	row := s.grid[y]
	return getCellInRow(row, x)
}

// setCell puts a stone of color c into cell (x, y)
// Cell (x, y) must be empty and valid
func (s *State) setCell(x, y byte, c color) {
	bits := uint64(c) << (x * 2)
	s.grid[y] |= bits
}

// getCellInRow returns color of a stone on index index in row row
func getCellInRow(row uint64, index byte) color {
	// Find the two bits that represent column with index index
	bits := ((3 << (index * 2)) & row) >> (index * 2)
	return getColorFromBits(bits)
}

func (s *State) clone() *State {
	return &State{s.size, s.grid, s.lastPlayer}
}

// IsCellValid returns true if a cell (x, y) is on the grid, and false otherwise
func (s *State) IsCellValid(x, y int) bool {
	return x >= 0 && x < int(s.size) && y >= 0 && y < int(s.size)
}

// IsCellEmpty returns true if a cell (x, y) is empty
func (s *State) IsCellEmpty(x, y byte) bool {
	return s.getColorOn(x, y) == None
}

// GetSuccessorState returns a state after Action a is performed
func (s *State) GetSuccessorState(a Action) *State {
	newState := s.clone()
	if a.c == s.lastPlayer {
		panic(fmt.Sprintf("Player cannot do two moves in a row! (last player: %s, current action: %s)", s.lastPlayer, a.c))
	}
	newState.lastPlayer = a.c
	newState.setCell(a.x, a.y, a.c)
	return newState
}

// GetPossibleActions returns a list of all possible actions from State s
func (s *State) GetPossibleActions() []Action {
	actions := make([]Action, 0, s.size*s.size)
	for rowIndex := byte(0); rowIndex < s.size; rowIndex++ {
		row := s.grid[rowIndex]
		for colIndex := byte(0); colIndex < s.size; colIndex++ {
			bits := row & 3 // Get last two bits of a row
			if getColorFromBits(bits) == None {
				actions = append(actions, Action{colIndex, rowIndex, s.lastPlayer.opponent()})
			}
			row = row >> 2
		}
	}
	return actions
}

// TODO:
// IsGoalState returns true if the game is decided (one player has a "virtual
// connection" and false otherwise)
func (s *State) IsGoalState() bool {
	return true
}
