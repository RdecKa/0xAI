package hexplayer

import (
	"errors"
	"fmt"

	"github.com/RdecKa/bachleor-thesis/1-mcts/mcts"
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

// MCTSplayer represents a computer player that uses only MCTS for selecting
// moves
type MCTSplayer struct {
	Color             hex.Color
	explorationFactor float64
	numIterations     int
	minBeforeExpand   uint
	mc                *mcts.MCTS
	state             *hex.State
	winPath           [][2]int
}

// CreateMCTSplayer creates a new player
func CreateMCTSplayer(c hex.Color, ef float64, ni int, mbe uint) *MCTSplayer {
	mp := MCTSplayer{c, ef, ni, mbe, nil, nil, nil}
	return &mp
}

// InitGame initializes the game
func (mp *MCTSplayer) InitGame(boardSize int) error {
	initState := hex.NewState(byte(boardSize))
	s := mcts.InitMCTS(*initState, mp.explorationFactor, mp.minBeforeExpand)
	mp.mc = s
	mp.state = initState
	return nil
}

// NextAction accepts opponent's last action and returns an action to be
// performed now.
func (mp *MCTSplayer) NextAction(prevAction *hex.Action) (*hex.Action, error) {
	// Update the state according to opponent's last move
	if prevAction != nil {
		s := mp.state.GetSuccessorState(prevAction).(hex.State)
		mp.state = &s
	}

	if len(mp.winPath) > 0 {
		fmt.Println("WIN PATH:")
		fmt.Println(mp.winPath)
		// TODO
		//return nil, nil
	}

	// Run MCTS
	mp.mc = mp.mc.ContinueMCTSFromChild(mp.state)
	if mp.mc == nil {
		return nil, errors.New("Cannot continue MCTS")
	}

	for i := 0; i < mp.numIterations; i++ {
		mp.mc.RunIteration()
	}

	// Get the best action
	bestState := mp.mc.GetBestRootChildState()
	bestAction := mp.state.GetTransitionAction(bestState).(*hex.Action)

	// Update mp.state
	s := bestState.(hex.State)
	mp.state = &s

	// Check if player has a virtual connection
	if exists, solution := mp.state.IsGoalState(false); exists {
		fmt.Println("Player has a virtual connection!")
		mp.winPath = solution.([][2]int)
	}

	return bestAction, nil
}

// EndGame doesn't do anything. The only reason for having it is that MCTSplayer
// must implement all functions of HexPlayer.
func (mp MCTSplayer) EndGame(lastAction *hex.Action, won bool) {}
