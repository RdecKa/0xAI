package hex

import (
	"fmt"
	"strconv"

	"github.com/RdecKa/bachleor-thesis/common/astarsearch"
	"github.com/RdecKa/bachleor-thesis/common/game"
)

// -----------------
// |     State     |
// -----------------

// State represents a state on a grid in a hex game
//	size is a length of the grid (size 11 means 11x11 grid)
//	grid is a list of rows in a grid, each row is represented as uint32. Each
//		cell in a row is stored with two bits:
//			00 - empty
//			01 - red
//			10 - blue
//			11 - undefined
//		Lowest two bits represent the cell with index 0. Because of using 32
//		bits for a row, maximal size of the grid is 16x16.
//	lastAction is the action that led to the current position on board
//	isInitialState tells whether the board is completely empty (at the beginning
//		of the game)
//
// A goal of the red player is to connect top-most and bottom-most row while a
// goal of the blue player is to connect left-most and right-most column
type State struct {
	size           byte
	grid           []uint32
	lastAction     *Action
	isInitialState bool
}

// NewState returns new State with a grid of given size and an invalid action as
// lastAction
func NewState(size byte, firstPlayer Color) *State {
	grid := make([]uint32, size)
	return &State{size, grid, NewAction(size, size, firstPlayer.Opponent()), true} // Opponent is set as last player, so firstPlayer starts
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

// GetSize returns size of the grid
func (s *State) GetSize() int {
	return int(s.size)
}

// GetCopyGrid returns a copy of the state's grid
func (s *State) GetCopyGrid() []uint32 {
	c := make([]uint32, len(s.grid))
	copy(c, s.grid)
	return c
}

// GetLastPlayer returns the player who made the last move
func (s State) GetLastPlayer() Color {
	return s.lastAction.c
}

// getColorOn returns the color of the stone in cell (x, y)
func (s *State) getColorOn(x, y byte) Color {
	row := s.grid[y]
	return getCellInRow(row, x)
}

// setCell puts a stone of color c into cell (x, y)
// Cell (x, y) must be empty and valid
func (s *State) setCell(x, y byte, c Color) {
	bits := uint32(c) << (x * 2)
	s.grid[y] |= bits
}

// getCellInRow returns color of a stone on index index in row row
func getCellInRow(row uint32, index byte) Color {
	// Find the two bits that represent column with index index
	bits := ((3 << (index * 2)) & row) >> (index * 2)
	return GetColorFromBits(bits)
}

func (s *State) clone() game.State {
	newGrid := make([]uint32, len(s.grid))
	for i, v := range s.grid {
		newGrid[i] = v
	}
	return State{s.size, newGrid, s.lastAction.clone(), s.isInitialState}
}

// IsCellValid returns true if a cell (x, y) is on the grid, and false otherwise
func (s *State) IsCellValid(x, y int) bool {
	return x >= 0 && x < int(s.size) && y >= 0 && y < int(s.size)
}

// IsCellEmpty returns true if a cell (x, y) is empty
func (s *State) IsCellEmpty(x, y byte) bool {
	return s.getColorOn(x, y) == None
}

// IsEndingCell returns true if a cell (x, y) is beyond player c's opposite edge
// - for red: below the board
// - for blue: to the right of the board
// Note: This is used to check whether a cell (x, y) is virtually connected to
// the opposite edge
func (s *State) IsEndingCell(x, y int, c Color) bool {
	return (c == Red && y >= int(s.size) && x >= 0) || (c == Blue && x >= int(s.size) && y >= 0)
}

// GetSuccessorState returns a state after Action a is performed
func (s State) GetSuccessorState(action game.Action) game.State {
	a := action.(*Action)
	if a.c == s.lastAction.c {
		panic(fmt.Sprintf("Player cannot do two moves in a row! (last player: %s, current action: %s)", s.lastAction.c, a))
	}
	if x, y := a.GetCoordinates(); s.getColorOn(byte(x), byte(y)) != None {
		panic(fmt.Sprintf("Cell (%d, %d) already occupied!", x, y))
	}
	newState := s.clone().(State)
	newState.lastAction.c = a.c
	newState.setCell(a.x, a.y, a.c)
	newState.isInitialState = false
	return newState
}

// GetPossibleActions returns a list of all possible actions from State s
func (s State) GetPossibleActions() []game.Action {
	actions := make([]game.Action, 0, s.size*s.size)
	playerColor := s.lastAction.c.Opponent()
	for rowIndex := byte(0); rowIndex < s.size; rowIndex++ {
		row := s.grid[rowIndex]
		for colIndex := byte(0); colIndex < s.size; colIndex++ {
			bits := row & 3 // Get last two bits of a row
			if GetColorFromBits(bits) == None {
				actions = append(actions, &Action{colIndex, rowIndex, playerColor})
			}
			row = row >> 2
		}
	}
	return actions
}

// IsGoalState returns true if the game is decided (the player who has just made
// a move has a (virtual) connection) and false otherwise
func (s State) IsGoalState(veryEnd bool) (bool, interface{}) {
	initialState := GetInitialState(&s)
	aStarSearch := astarsearch.InitSearch(&initialState)
	solutionExists, solution := aStarSearch.Search(veryEnd)
	return solutionExists, solution
}

// EvaluateGoalState returns 1.0 because the player who makes the last action
// (action that leads to the goal state) wins
func (s State) EvaluateGoalState() float64 {
	return 1.0
}

// Same returns true if states s and s2 represent the same state on the board.
func (s State) Same(sg game.State) bool {
	s2 := sg.(*State)
	if s.size != s2.size {
		return false
	}
	if s.lastAction.c != s2.lastAction.c {
		return false
	}
	for i := byte(0); i < s.size; i++ {
		if s.grid[i] != s2.grid[i] {
			return false
		}
	}
	return true
}

// GetNumOfStones returns number of red stones, blue stones and empty cells (in
// that order)
func (s State) GetNumOfStones() (int, int, int) {
	red, blue, empty := 0, 0, 0
	for _, row := range s.grid {
		r := row
		for colIndex := byte(0); colIndex < s.size; colIndex++ {
			c := GetColorFromBits(r & 3)
			switch c {
			case Red:
				red++
			case Blue:
				blue++
			default:
				empty++
			}
			r = r >> 2
		}
	}
	return red, blue, empty
}

// GetTransitionAction returns an action that leads from State s to State sg.
func (s State) GetTransitionAction(sg game.State) game.Action {
	s2 := sg.(State)
	for r := 0; r < s.GetSize(); r++ {
		if s.grid[r] != s2.grid[r] {
			diff := s2.grid[r] - s.grid[r]
			c := 0
			for diff > 3 {
				c++
				diff = diff >> 2
			}
			return NewAction(byte(c), byte(r), GetColorFromBits(diff))
		}
	}
	return nil
}

// GetMapKey generates a key to be used in a hash map
func (s State) GetMapKey() string {
	st := ""
	for _, row := range s.grid {
		st += strconv.Itoa(int(row)) + ","
	}
	return st
}
