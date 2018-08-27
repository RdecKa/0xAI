// Package hexplayer - In this file are implemented functions that can be called by any coputer
// player when the game is decided.
package hexplayer

import (
	"errors"
	"fmt"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

// findSafeCells returns all the bridges on the winning path.
func findSafeCells(winPath [][2]int, boardSize int, playerColor hex.Color) [][2]cell {
	safeCells := make([][2]cell, 0)
	if playerColor == hex.Red {
		if winPath[0][1] < boardSize-1 {
			// Bridge to the bottom (bottom row does not have a stone yet)
			x := winPath[0][0]
			y := winPath[0][1]
			safeCells = append(safeCells, returnTwoSafeCellsBetween(x, y, x-1, y+2))
		}
		if winPath[len(winPath)-2][1] > 0 {
			// Bridge to the top
			x := winPath[len(winPath)-2][0]
			y := winPath[len(winPath)-2][1]
			safeCells = append(safeCells, returnTwoSafeCellsBetween(x, y, x+1, y-2))
		}
	} else if playerColor == hex.Blue {
		if winPath[0][0] < boardSize-1 {
			// Bridge to the right
			x := winPath[0][0]
			y := winPath[0][1]
			safeCells = append(safeCells, returnTwoSafeCellsBetween(x, y, x+2, y-1))
		}
		if winPath[len(winPath)-2][0] > 0 {
			// Bridge to the left
			x := winPath[len(winPath)-2][0]
			y := winPath[len(winPath)-2][1]
			safeCells = append(safeCells, returnTwoSafeCellsBetween(x, y, x-2, y+1))
		}
	}
	for i := range winPath[:len(winPath)-2] {
		if !(directlyConnected(winPath[i], winPath[i+1])) {
			safeCells = append(safeCells, returnTwoSafeCellsBetween(winPath[i][0], winPath[i][1], winPath[i+1][0], winPath[i+1][1]))
		}
	}
	return safeCells
}

// returnTwoSafeCellsBetween returns two cells between (x1, y1) and (x2, y2).
// (x1, y1) -- (x2, y2) must be a bridge, this function does not check that.
func returnTwoSafeCellsBetween(x1, y1, x2, y2 int) [2]cell {
	var nx1, nx2, ny1, ny2 int
	if (x1-x2)%2 == 0 {
		nx1 = (x1 + x2) / 2
		nx2 = nx1
	} else {
		nx1 = (x1 + x2) / 2
		nx2 = nx1 + 1
	}

	if (y1-y2)%2 == 0 {
		ny1 = (y1 + y2) / 2
		ny2 = ny1
	} else {
		ny2 = (y1 + y2) / 2
		ny1 = ny2 + 1
	}

	return [2]cell{cell{nx1, ny1}, cell{nx2, ny2}}
}

// directlyConnected checks whether two cells are directly connected (they share
// one side)
func directlyConnected(c1, c2 [2]int) bool {
	x1, y1 := c1[0], c1[1]
	x2, y2 := c2[0], c2[1]
	if x1 == x2 && abs(y1-y2) == 1 {
		return true
	}
	if y1 == y2 && abs(x1-x2) == 1 {
		return true
	}
	if abs(x1-x2) == 1 && abs(y1-y2) == 1 && (x1-x2)+(y1-y2) == 0 {
		return true
	}
	return false
}

func abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

// getAttackedBridge returns two integers:
//	- index of the bridge that was attacked with Action prevAction (-1 if none of the bridges is attacked)
//	- index of the cell in the attacked bridge that is still empty (-1 if none of the bridges is attacked)
func getAttackedBridge(prevAction *hex.Action, safeWinCells [][2]cell) (int, int) {
	for i, c := range safeWinCells {
		x, y := prevAction.GetCoordinates()
		if x == c[0].x && y == c[0].y {
			return i, 1
		}
		if x == c[1].x && y == c[1].y {
			return i, 0
		}
	}
	return -1, -1
}

// doNotLoseHope is called when the game is lost but player still want to attack
// the opponent and wait for his/her mistake.
func doNotLoseHope(state *hex.State, playerColor hex.Color) (*hex.Action, error) {
	fmt.Println("A player doesn't want to give up!")
	exists, solution := state.IsGoalState(false)
	if !exists {
		return nil, errors.New("Game lost but solution does not exist")
	}
	winPath := solution.([][2]int)

	// Find opponent's safe cells, attack one
	safeCells := findSafeCells(winPath, state.GetSize(), playerColor.Opponent())
	a := hex.NewAction(byte(safeCells[0][0].x), byte(safeCells[0][0].y), playerColor)

	return a, nil
}

// getActionIfWinningPathExists checks whether the winning path already exists.
// If it does (indicated by bool return value), it returns an action that either
// defends the attacked bridge (if there is an attacked bridge) or connects one
// of the remaining open bridges. In the latter case, the updated safeWinCells
// is also returned.
func getActionIfWinningPathExists(lastOpponentAction *hex.Action, safeWinCells [][2]cell, playerColor hex.Color) (*hex.Action, [][2]cell, bool) {
	if len(safeWinCells) <= 0 {
		// Winning path does not exist (yet)
		return nil, nil, false
	}

	var ec cell
	if bridge, emptyCellIndex := getAttackedBridge(lastOpponentAction, safeWinCells); bridge > -1 {
		// Opponent has attacked one of the bridges
		ec = safeWinCells[bridge][emptyCellIndex]
		safeWinCells = append(safeWinCells[:bridge], safeWinCells[bridge+1:]...)
	} else {
		// Select the first cell in the first bridge (doesn't really matter, which one)
		ec = safeWinCells[0][0]
		safeWinCells = safeWinCells[1:]
	}

	action := hex.NewAction(byte(ec.x), byte(ec.y), playerColor)

	return action, safeWinCells, true
}
