package hex

import (
	"fmt"

	"github.com/RdecKa/common/astarsearch"
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

// Cells in a rectangular grid (x, y) that are virtually connected to cell
// (0, 0), if and only if the two cells between (0, 0) and (x, y) are empty
var virtualConnections = [6][]int{
	[]int{-1, -1},
	[]int{1, -2},
	[]int{2, -1},
	[]int{1, 1},
	[]int{-1, 2},
	[]int{-2, 1},
}

type searchState struct {
	x, y      int    // Current position; x = -1 means a field left of the grid, y = -1 means a field above the grid
	c         color  // The color that a solution is searched for
	gameState *State // State in a game where solution is searched for
}

// GetInitialState returns inital state of the search
func GetInitialState(gameState *State) searchState {
	return searchState{-1, -1, gameState.lastPlayer, gameState}
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
		xp, yp := &a, &b
		if s.c == Blue {
			xp, yp = yp, xp
		}

		// Add cells in the first row/column (directly connected to the edge)
		for ; a < int(s.gameState.size); a++ {
			if s.gameState.getColorOn(byte(*xp), byte(*yp)) == s.c {
				successors = append(successors, searchState{*xp, *yp, s.c, s.gameState})
			}
		}

		// Add cells in the second row/column (those which are virtually
		// connected to the player's first edge)
		b = 1
		xDiffFirst, yDiffFirst, xDiffSecond, yDiffSecond := 0, -1, 1, -1
		if s.c == Blue {
			xDiffFirst, yDiffFirst, xDiffSecond, yDiffSecond = yDiffFirst, xDiffFirst, yDiffSecond, xDiffSecond
		}
		for a = 0; a < int(s.gameState.size)-1; a++ {
			if s.gameState.getColorOn(byte(*xp), byte(*yp)) == s.c &&
				s.gameState.getColorOn(byte(*xp+xDiffFirst), byte(*yp+yDiffFirst)) == None &&
				s.gameState.getColorOn(byte(*xp+xDiffSecond), byte(*yp+yDiffSecond)) == None {
				successors = append(successors, searchState{*xp, *yp, s.c, s.gameState})
			}
		}
	} else {
		// Add direct neighbours
		for _, n := range neighbours {
			x, y := s.x+n[0], s.y+n[1]
			if s.gameState.IsCellValid(x, y) && s.gameState.getColorOn(byte(x), byte(y)) == s.c {
				successors = append(successors, searchState{x, y, s.c, s.gameState})
			}
		}
		// Add virtual connections
		for _, v := range virtualConnections {
			// Check if the two cells between current cell (s.x, s.y) and (x, y)
			// are empty
			var x1, x2, y1, y2 int // Relative coordinates of these two cells
			if v[0]%2 == 0 {       // difference in x coordinates is 2 or -2
				x1, x2 = v[0]/2, v[0]/2
			} else { // difference in x coordinates is 1 or -1
				x1, x2 = v[0], 0
			}

			if v[1]%2 == 0 { // difference in y coordinates is 2 or -2
				y1, y2 = v[1]/2, v[1]/2
			} else {
				y1, y2 = 0, v[1] // difference in y coordinates is 1 or -1
			}

			// Change relative coordinates to absolute coordinates
			x1, x2, y1, y2 = s.x+x1, s.x+x2, s.y+y1, s.y+y2

			if s.gameState.IsCellValid(x1, y1) && s.gameState.getColorOn(byte(x1), byte(y1)) != None ||
				s.gameState.IsCellValid(x2, y2) && s.gameState.getColorOn(byte(x2), byte(y2)) != None {
				continue // At least one of these cells is not empty
			}
			// Both cells in between are empty

			x, y := s.x+v[0], s.y+v[1]

			if s.gameState.IsEndingCell(x, y, s.c) ||
				s.gameState.IsCellValid(x, y) && s.gameState.getColorOn(byte(x), byte(y)) == s.c {
				successors = append(successors, searchState{x, y, s.c, s.gameState})
			}
		}
	}
	return successors
}

func (s searchState) String() string {
	return fmt.Sprintf("x: %d, y: %d, c: %s, state:\n%s", s.x, s.y, s.c, s.gameState.String())
}
