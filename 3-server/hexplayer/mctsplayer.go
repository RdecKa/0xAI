package hexplayer

import (
	"math/rand"

	"github.com/RdecKa/bachleor-thesis/1-mcts/mcts"
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

type MCTSplayer struct {
	Color             hex.Color
	explorationFactor float64
	numIterations     int
	minBeforeExpand   uint
	mc                *mcts.MCTS
	state             *hex.State
}

func CreateMCTSplayer(c hex.Color, ef float64, ni int, mbe uint) *MCTSplayer {
	mp := MCTSplayer{c, ef, ni, mbe, nil, nil}
	return &mp
}

func (mp *MCTSplayer) InitGame(boardSize int) error {
	initState := hex.NewState(byte(boardSize))
	s := mcts.InitMCTS(*initState, mp.explorationFactor, mp.minBeforeExpand)
	mp.mc = s
	mp.state = initState
	return nil
}

func (mp *MCTSplayer) NextAction(prevAction *hex.Action) (*hex.Action, error) {
	if prevAction != nil {
		s := mp.state.GetSuccessorState(prevAction).(hex.State)
		mp.state = &s
	}

	/*mp.mc = mp.mc.ContinueMCTSFromChild(mp.state)
	if mp.mc == nil {
		return nil, errors.New("Cannot continue MCTS")
	}

	for i := 0; i < mp.numIterations; i++ {
		mp.mc.RunIteration()
	}

	// Find a way to get the best action

	// Update mp.state

	return nil, nil*/

	// For now, just random
	as := mp.state.GetPossibleActions()
	r := as[rand.Intn(len(as))].(*hex.Action)
	s := mp.state.GetSuccessorState(r).(hex.State)
	mp.state = &s
	return r, nil
}

func (mp MCTSplayer) EndGame(lastAction *hex.Action, won bool) {
}
