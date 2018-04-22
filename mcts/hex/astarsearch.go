package hex

import (
	"fmt"

	"github.com/RdecKa/mcts/astarsearch"
)

// Cells in a rectangular grid that are neighbours in a hexagonal grid
var neighbours = [6][]int{
	[]int{0, -1},
	[]int{1, -1},
	[]int{-1, 0},
	[]int{1, 0},
	[]int{-1, 1},
	[]int{0, 1},
}

type searchState struct {
	x, y      int    // Current position; x = -1 means a field left of the grid, y = -1 means a field above the grid
	c         color  // The color that a solution is searched for
	gameState *State // State in a game where solution is searched for
}

// GetInitialState returns inital state of the search
func GetInitialState(c color, gameState *State) searchState {
	return searchState{-1, -1, c, gameState}
}

func (s searchState) IsGoalState() bool {
	size := int(s.gameState.GetSize())
	return (s.c == Red && s.y >= size-1) || (s.c == Blue && s.x >= size-1)
}

// getEstimateToReachGoal returns number of cells between current state and
// final row/column
func (s searchState) GetEstimateToReachGoal() int {
	switch s.c {
	case Red:
		return s.gameState.GetSize() - 1 - s.y
	case Blue:
		return s.gameState.GetSize() - 1 - s.x
	}
	panic(fmt.Sprintf("Unknown color %d!", s.c))
}

func (s searchState) GetSuccessorStates() []astarsearch.State {
	successors := make([]astarsearch.State, 0, 6)

	if (s.x == -1 && s.c == Blue) || (s.y == -1 && s.c == Red) {
		// Beginning of the search
		a, b := 0, 0
		ap, bp := &a, &b
		if s.c == Blue {
			ap, bp = bp, ap
		}

		for ; a < int(s.gameState.size); a++ {
			if s.gameState.getColorOn(byte(*ap), byte(*bp)) == s.c {
				newState := searchState{*ap, *bp, s.c, s.gameState}
				successors = append(successors, newState)
			}
		}
	} else {
		for _, n := range neighbours {
			x, y := s.x+n[0], s.y+n[1]
			if s.gameState.IsCellValid(x, y) && s.gameState.getColorOn(byte(x), byte(y)) == s.c {
				newState := searchState{x, y, s.c, s.gameState}
				successors = append(successors, newState)
			}
		}
	}
	return successors
}
