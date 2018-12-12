package hexplayer

import (
	"errors"
	"fmt"
	"time"

	"github.com/RdecKa/0xAI/1-mcts/mcts"
	"github.com/RdecKa/0xAI/common/game/hex"
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
	allowResignation   bool
}

// CreateMCTSplayer creates a new player
func CreateMCTSplayer(c hex.Color, ef float64, t time.Duration, mbe uint, ar bool) *MCTSplayer {
	mp := MCTSplayer{c, ef, t, mbe, nil, nil, nil, 0, nil, ar}
	return &mp
}

// InitGame initializes the game
func (mp *MCTSplayer) InitGame(boardSize int, firstPlayer hex.Color) error {
	initState := hex.NewState(byte(boardSize), firstPlayer)
	mp.mc = mcts.InitMCTS(*initState, mp.explorationFactor, mp.minBeforeExpand)
	mp.state = initState
	mp.safeWinCells = nil
	mp.lastOpponentAction = nil
	return nil
}

func (mp *MCTSplayer) initGameFromState(initState *hex.State, lastOpponentAction *hex.Action) error {
	mp.mc = mcts.InitMCTS(initState, mp.explorationFactor, mp.minBeforeExpand)
	mp.state = initState
	mp.lastOpponentAction = lastOpponentAction

	if exists, solution := initState.IsGoalState(false); exists {
		winPath := solution.([][2]int)
		mp.safeWinCells = findSafeCells(winPath, initState.GetSize(), mp.Color)
	} else {
		mp.safeWinCells = nil
	}

	return nil
}

// PrevAction accepts opponent's last action
func (mp *MCTSplayer) PrevAction(prevAction *hex.Action) {
	// Update the state according to opponent's last move
	if prevAction != nil {
		mp.updatePlayerState(prevAction)
		mp.lastOpponentAction = prevAction
	}
}

// NextAction returns an action to be performed. It returns nil when the player
// decides to resign.
func (mp *MCTSplayer) NextAction() (*hex.Action, error) {
	// Check if the player has already won (has a virtual connection)
	if a, swc, ok := getActionIfWinningPathExists(mp.lastOpponentAction, mp.safeWinCells, mp.Color); ok {
		mp.updatePlayerState(a)
		mp.safeWinCells = swc
		return a, nil
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
			mp.mc.RunIteration(true)
		}
	}

	// Get the best action
	bestState := mp.mc.GetBestRootChildState()
	if bestState == nil {
		// Game lost
		if !mp.allowResignation {
			// Continue playing and hope for opponent's mistake
			a, err := doNotLoseHope(mp.state, mp.Color)
			if err != nil {
				return nil, err
			}

			// Update state
			mp.updatePlayerState(a)

			return a, nil
		}
		// Resign
		return nil, nil
	}
	bestAction := bestState.(hex.State).GetLastAction()

	// Update mp.state
	s := bestState.(hex.State)
	mp.state = &s

	// Check if player has a virtual connection
	if exists, solution := mp.state.IsGoalState(false); exists {
		fmt.Println("MCTS Player has a virtual connection!")
		winPath := solution.([][2]int)
		mp.safeWinCells = findSafeCells(winPath, mp.state.GetSize(), mp.Color)
	}

	return bestAction, nil
}

// updatePlayerState updates the game state of the player
func (mp *MCTSplayer) updatePlayerState(a *hex.Action) {
	s := mp.state.GetSuccessorState(a).(hex.State)
	mp.state = &s
}

// EndGame accepts the result of the game
func (mp *MCTSplayer) EndGame(lastAction *hex.Action, won bool) {
	if won {
		mp.numWin++
	}
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
