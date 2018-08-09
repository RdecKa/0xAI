package ab

import (
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

const bigNumber = 10000.0

// AlphaBeta runs search with AB pruning to select the next action to be taken
func AlphaBeta(state *hex.State) *hex.Action {
	gridChan, stopChan, resultChan := hex.CreatePatChecker()
	defer func() { stopChan <- struct{}{} }()

	_, a := alphaBeta(2, state, -bigNumber, bigNumber, gridChan, resultChan)

	if a == nil {
		possibleActions := state.GetPossibleActions()
		return possibleActions[0].(*hex.Action)
	}

	return a
}

func alphaBeta(depth int, state *hex.State, alpha, beta float64, gridChan chan []uint64, resultChan chan [2][]int) (float64, *hex.Action) {
	goal, _ := state.IsGoalState(false)
	if goal || depth <= 0 {
		return eval(state, gridChan, resultChan), nil
	}
	bestValue := bigNumber
	var bestState hex.State

	possibleActions := state.GetPossibleActions()
	for _, a := range possibleActions {
		successor := state.GetSuccessorState(a).(hex.State)
		value, _ := alphaBeta(depth-1, &successor, -beta, -alpha, gridChan, resultChan)
		value = -value

		bestValue = value
		bestState = successor

		if bestValue > alpha {
			alpha = bestValue
			if alpha >= beta {
				break
			}
		}
	}

	return bestValue, state.GetTransitionAction(bestState).(*hex.Action)
}

func eval(state *hex.State, gridChan chan []uint64, resultChan chan [2][]int) float64 {
	gridChan <- state.GetCopyGrid()
	red, blue, _ := state.GetNumOfStones()
	patCount := <-resultChan
	sample := Sample{
		num_stones:    red + blue,
		occ_red_rows:  patCount[0][len(patCount)-2],
		occ_red_cols:  patCount[0][len(patCount)-1],
		occ_blue_rows: patCount[1][len(patCount)-2],
		occ_blue_cols: patCount[1][len(patCount)-1],
		red_p1:        patCount[0][0],
		blue_p1:       patCount[1][0],
		red_p2:        patCount[0][1],
		blue_p2:       patCount[1][1],
	}
	return sample.getEstimatedValue()
}
