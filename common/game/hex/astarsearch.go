package hex

import (
	"fmt"

	"github.com/RdecKa/bachleor-thesis/common/astarsearch"
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

// Assumption: c and gameState are the same in all searchStates that appear in
// one run of A* search.
var currentGameState *State // State in a game where solution is searched for
var currentColor Color      // The color that a solution is searched for

type searchState struct {
	x, y            int          // Current position; x = -1 means a field left of the grid, y = -1 means a field above the grid
	prevSearchState *searchState // The preceeding search state
}

// GetInitialState returns inital state of the search
func GetInitialState(gameState *State) searchState {
	currentGameState = gameState
	currentColor = gameState.lastPlayer
	return searchState{-1, -1, nil}
}

// GetClean returns an array that is unique for each searchState.
// Can be used as a key in a map.
func (s searchState) GetClean() interface{} {
	return [2]int{s.x, s.y}
}

func (s searchState) IsGoalState() (bool, interface{}) {
	size := int(currentGameState.GetSize())
	isGoal := (currentColor == Red && s.y >= size-1) || (currentColor == Blue && s.x >= size-1)
	if isGoal {
		return true, s.GetWinPath()
	}
	return false, nil
}

// GetEstimateToReachGoal returns number of cells between current state and
// final row/column
func (s searchState) GetEstimateToReachGoal() int {
	switch currentColor {
	case Red:
		return currentGameState.GetSize() - 1 - s.y
	case Blue:
		return currentGameState.GetSize() - 1 - s.x
	}
	panic(fmt.Sprintf("Unknown color %d!", currentColor))
}

// GetSuccessorStates returns all possible successors of the searchState s.
// If veryEnd == true, then the game has actually ended
// Else if veryEnd == false, then the game has theoretically ended
func (s searchState) GetSuccessorStates(veryEnd bool) []astarsearch.State {
	successors := make([]astarsearch.State, 0, 6)

	if (s.x == -1 && currentColor == Blue) || (s.y == -1 && currentColor == Red) {
		// Beginning of the search
		a, b := 0, 0
		xp, yp := &a, &b
		if currentColor == Blue {
			xp, yp = yp, xp
		}

		// Add cells in the first row/column (directly connected to the edge)
		for ; a < int(currentGameState.size); a++ {
			if currentGameState.getColorOn(byte(*xp), byte(*yp)) == currentColor {
				successors = append(successors, searchState{*xp, *yp, &s})
			}
		}

		if !veryEnd {
			// Add cells in the second row/column (those which are virtually
			// connected to the player's first edge)
			b = 1
			xDiffFirst, yDiffFirst, xDiffSecond, yDiffSecond := 0, -1, 1, -1
			if currentColor == Blue {
				xDiffFirst, yDiffFirst, xDiffSecond, yDiffSecond = yDiffFirst, xDiffFirst, yDiffSecond, xDiffSecond
			}
			for a = 0; a < int(currentGameState.size)-1; a++ {
				if currentGameState.getColorOn(byte(*xp), byte(*yp)) == currentColor &&
					currentGameState.getColorOn(byte(*xp+xDiffFirst), byte(*yp+yDiffFirst)) == None &&
					currentGameState.getColorOn(byte(*xp+xDiffSecond), byte(*yp+yDiffSecond)) == None {
					successors = append(successors, searchState{*xp, *yp, &s})
				}
			}
		}
	} else {
		// Add direct neighbours
		for _, n := range neighbours {
			x, y := s.x+n[0], s.y+n[1]
			if currentGameState.IsCellValid(x, y) && currentGameState.getColorOn(byte(x), byte(y)) == currentColor {
				successors = append(successors, searchState{x, y, &s})
			}
		}
		if !veryEnd {
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

				if currentGameState.IsCellValid(x1, y1) && currentGameState.getColorOn(byte(x1), byte(y1)) != None ||
					currentGameState.IsCellValid(x2, y2) && currentGameState.getColorOn(byte(x2), byte(y2)) != None {
					continue // At least one of these cells is not empty
				}
				// Both cells in between are empty

				x, y := s.x+v[0], s.y+v[1]

				if currentGameState.IsEndingCell(x, y, currentColor) ||
					currentGameState.IsCellValid(x, y) && currentGameState.getColorOn(byte(x), byte(y)) == currentColor {
					successors = append(successors, searchState{x, y, &s})
				}
			}
		}
	}
	return successors
}

func (s searchState) String() string {
	return fmt.Sprintf("x: %d, y: %d, c: %s, state:\n%s", s.x, s.y, currentColor, currentGameState.String())
}

// GetWinPath returns coordinates of the cells in the winning (virtual)
// connection.
func (s *searchState) GetWinPath() [][2]int {
	path := make([][2]int, 0, currentGameState.GetSize())

	cs := s
	for cs != nil {
		path = append(path, [2]int{cs.x, cs.y})
		cs = cs.prevSearchState
	}

	return path
}
