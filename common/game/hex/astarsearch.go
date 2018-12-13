package hex

import (
	"fmt"

	"github.com/RdecKa/0xAI/common/astarsearch"
)

// Cells in a rectangular grid that are neighbours in a hexagonal grid
// DO NOT MIX THE ORDER
var neighbours = [6][]int{
	[]int{0, -1},
	[]int{1, -1},
	[]int{1, 0},
	[]int{0, 1},
	[]int{-1, 1},
	[]int{-1, 0},
}

// Cells in a rectangular grid (x, y) that are virtually connected to cell
// (0, 0), if and only if the two cells between (0, 0) and (x, y) are empty
// DO NOT MIX THE ORDER
var virtualConnections = [6][]int{
	[]int{-1, -1},
	[]int{1, -2},
	[]int{2, -1},
	[]int{1, 1},
	[]int{-1, 2},
	[]int{-2, 1},
}

type searchState struct {
	x, y            int          // Current position; x = -1 means a field left of the grid, y = -1 means a field above the grid
	c               Color        // The color that a solution is searched for
	gameState       *State       // State in a game where solution is searched for
	prevSearchState *searchState // The preceeding search state
}

// GetInitialState returns inital state of the search
func GetInitialState(gameState *State) searchState {
	return searchState{-1, -1, gameState.lastAction.c, gameState, nil}
}

// GetClean returns an array that is unique for each searchState.
// Can be used as a key in a map.
func (s searchState) GetClean() interface{} {
	return [2]int{s.x, s.y}
}

func (s searchState) IsGoalState() (bool, interface{}) {
	size := int(s.gameState.GetSize())
	isGoal := (s.c == Red && s.y >= size-1) || (s.c == Blue && s.x >= size-1)
	if isGoal {
		return true, s.GetWinPath()
	}
	return false, nil
}

// GetEstimateToReachGoal returns number of cells between current state and
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

// GetSuccessorStates returns all possible successors of the searchState s.
// If veryEnd == true, then the game has actually ended
// Else if veryEnd == false, then the game has theoretically ended
func (s searchState) GetSuccessorStates(veryEnd bool) []astarsearch.State {
	successors := make([]astarsearch.State, 0)

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
				successors = append(successors, searchState{*xp, *yp, s.c, s.gameState, &s})
			}
		}

		if !veryEnd {
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
					successors = append(successors, searchState{*xp, *yp, s.c, s.gameState, &s})
				}
			}
		}
	} else {
		// Add direct neighbours
		directNeighboursEmpty := make([]bool, 6)
		for in, n := range neighbours {
			x, y := s.x+n[0], s.y+n[1]
			if s.gameState.IsCellValid(x, y) {
				c := s.gameState.getColorOn(byte(x), byte(y))
				if c == s.c {
					successors = append(successors, searchState{x, y, s.c, s.gameState, &s})
				} else if c == None {
					directNeighboursEmpty[in] = true
				}
			}
		}
		if !veryEnd {
			// Add virtual connections
			for ivc, vc := range virtualConnections {
				if !directNeighboursEmpty[(ivc+5)%6] || !directNeighboursEmpty[ivc] {
					continue // At least one of these cells is not empty
				}
				// Both cells in between are empty

				x, y := s.x+vc[0], s.y+vc[1]

				if s.gameState.IsEndingCell(x, y, s.c) ||
					s.gameState.IsCellValid(x, y) &&
						s.gameState.getColorOn(byte(x), byte(y)) == s.c &&
						!s.overlapsWithIncomingVirtualConnection(x, y) {
					successors = append(successors, searchState{x, y, s.c, s.gameState, &s})
				}
			}
		}
	}
	return successors
}

func (s searchState) overlapsWithIncomingVirtualConnection(x, y int) bool {
	curCoords := [2]int{s.x, s.y}
	newCoords := [2]int{x, y}
	if DirectlyConnected(newCoords, curCoords) {
		return false
	}

	prvCoords := [2]int{s.prevSearchState.x, s.prevSearchState.y}
	if DirectlyConnected(curCoords, prvCoords) {
		return false
	}

	lastConn := GetTwoCellsBewteen(curCoords, prvCoords)
	nextConn := GetTwoCellsBewteen(curCoords, newCoords)

	if lastConn[0] == nextConn[0] || lastConn[0] == nextConn[1] ||
		lastConn[1] == nextConn[0] || lastConn[1] == nextConn[1] {
		return true
	}

	return false
}

func (s searchState) String() string {
	return fmt.Sprintf("x: %d, y: %d, c: %s, state:\n%s", s.x, s.y, s.c, s.gameState.String())
}

// GetWinPath returns coordinates of the cells in the winning (virtual)
// connection.
func (s *searchState) GetWinPath() [][2]int {
	path := make([][2]int, 0, s.gameState.GetSize())

	cs := s
	for cs != nil {
		path = append(path, [2]int{cs.x, cs.y})
		cs = cs.prevSearchState
	}

	return path
}
