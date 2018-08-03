package hexplayer

import (
	"errors"
	"fmt"
	"time"

	"github.com/RdecKa/bachleor-thesis/1-mcts/mcts"
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

// MCTSplayer represents a computer player that uses only MCTS for selecting
// moves
type MCTSplayer struct {
	Color              hex.Color
	explorationFactor  float64
	timeToRun          time.Duration
	minBeforeExpand    uint
	mc                 *mcts.MCTS
	state              *hex.State
	safeWinCells       [][2]cell
	numWin             int
	lastOpponentAction *hex.Action
}

type cell struct {
	x, y int
}

// CreateMCTSplayer creates a new player
func CreateMCTSplayer(c hex.Color, ef float64, t time.Duration, mbe uint) *MCTSplayer {
	mp := MCTSplayer{c, ef, t, mbe, nil, nil, nil, 0, nil}
	return &mp
}

// InitGame initializes the game
func (mp *MCTSplayer) InitGame(boardSize int, firstPlayer hex.Color) error {
	initState := hex.NewState(byte(boardSize), firstPlayer)
	s := mcts.InitMCTS(*initState, mp.explorationFactor, mp.minBeforeExpand)
	mp.mc = s
	mp.state = initState
	return nil
}

// PrevAction accepts opponent's last action
func (mp *MCTSplayer) PrevAction(prevAction *hex.Action) {
	// Update the state according to opponent's last move
	if prevAction != nil {
		s := mp.state.GetSuccessorState(prevAction).(hex.State)
		mp.state = &s
		mp.lastOpponentAction = prevAction
	}
}

// NextAction returns an action to be performed. It returns nil when it decides
// to resign.
func (mp *MCTSplayer) NextAction() (*hex.Action, error) {
	if len(mp.safeWinCells) > 0 {
		var ec cell
		if bridge, emptyCellIndex := mp.getAttackedBridge(mp.lastOpponentAction); bridge > -1 {
			// Opponent has attacked one of the bridges
			ec = mp.safeWinCells[bridge][emptyCellIndex]
			mp.safeWinCells = append(mp.safeWinCells[:bridge], mp.safeWinCells[bridge+1:]...)
		} else {
			// Select the first cell in the first bridge (doesn't really matter, which one)
			ec = mp.safeWinCells[0][0]
			mp.safeWinCells = mp.safeWinCells[1:]
		}
		action := hex.NewAction(byte(ec.x), byte(ec.y), mp.Color)
		s := mp.state.GetSuccessorState(action).(hex.State)
		mp.state = &s
		return hex.NewAction(byte(ec.x), byte(ec.y), mp.Color), nil
	}

	// Run MCTS
	mp.mc = mp.mc.ContinueMCTSFromChild(mp.state)
	if mp.mc == nil {
		return nil, errors.New("Cannot continue MCTS")
	}

	timer := time.NewTimer(mp.timeToRun)

	for timeOut := false; !timeOut; {
		select {
		case <-timer.C:
			timeOut = true
		default:
			mp.mc.RunIteration()
		}
	}

	// Get the best action
	bestState := mp.mc.GetBestRootChildState()
	if bestState == nil {
		// Game lost, resign
		return nil, nil
	}
	bestAction := mp.state.GetTransitionAction(bestState).(*hex.Action)

	// Update mp.state
	s := bestState.(hex.State)
	mp.state = &s

	// Check if player has a virtual connection
	if exists, solution := mp.state.IsGoalState(false); exists {
		fmt.Println("Player has a virtual connection!")
		winPath := solution.([][2]int)
		mp.findSafeCells(winPath)
	}

	return bestAction, nil
}

// EndGame accepts the result of the game
func (mp *MCTSplayer) EndGame(lastAction *hex.Action, won bool) {
	if won {
		mp.numWin++
	}
}

// findSafeCells saves all the bridges on the winning path.
func (mp *MCTSplayer) findSafeCells(winPath [][2]int) {
	if mp.Color == hex.Red {
		if winPath[0][1] < mp.state.GetSize()-1 {
			// Bridge to the bottom (bottom row does not have a stone yet)
			x := winPath[0][0]
			y := winPath[0][1]
			mp.addTwoSafeCellsBetween(x, y, x-1, y+2)
		}
		if winPath[len(winPath)-2][1] > 0 {
			// Bridge to the top
			x := winPath[len(winPath)-2][0]
			y := winPath[len(winPath)-2][1]
			mp.addTwoSafeCellsBetween(x, y, x+1, y-2)
		}
	} else if mp.Color == hex.Blue {
		if winPath[0][0] < mp.state.GetSize()-1 {
			// Bridge to the right
			x := winPath[0][0]
			y := winPath[0][1]
			mp.addTwoSafeCellsBetween(x, y, x+2, y-1)
		}
		if winPath[len(winPath)-2][0] > 0 {
			// Bridge to the left
			x := winPath[len(winPath)-2][0]
			y := winPath[len(winPath)-2][1]
			mp.addTwoSafeCellsBetween(x, y, x-2, y+1)
		}
	}
	for i := range winPath[:len(winPath)-2] {
		if !(directlyConnected(winPath[i], winPath[i+1])) {
			mp.addTwoSafeCellsBetween(winPath[i][0], winPath[i][1], winPath[i+1][0], winPath[i+1][1])
		}
	}
}

// addTwoSafeCellsBetween adds two cells between (x1, y1) and (x2, y2) to the
// list of safe cells. (x1, y1) -- (x2, y2) must be a bridge, this function does
// not check that.
func (mp *MCTSplayer) addTwoSafeCellsBetween(x1, y1, x2, y2 int) {
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

	mp.safeWinCells = append(mp.safeWinCells, [2]cell{cell{nx1, ny1}, cell{nx2, ny2}})
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
func (mp *MCTSplayer) getAttackedBridge(prevAction *hex.Action) (int, int) {
	for i, c := range mp.safeWinCells {
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

// SwitchColor switches the color of the player
func (mp *MCTSplayer) SwitchColor() {
	mp.Color = mp.Color.Opponent()
}

// GetColor returns the color of the player
func (mp MCTSplayer) GetColor() hex.Color {
	return mp.Color
}

// GetNumberOfWins returns the number of wins for this player
func (mp MCTSplayer) GetNumberOfWins() int {
	return mp.numWin
}

// GetType returns the type of the player
func (mp MCTSplayer) GetType() PlayerType {
	return MctsType
}
