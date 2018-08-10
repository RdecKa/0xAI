package ab

import (
	"fmt"
	"log"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

const maxValue = 10000.0
const abInit = 1000.0
const won = 500.0

// AlphaBeta runs search with AB pruning to select the next action to be taken
func AlphaBeta(state *hex.State) *hex.Action {
	gridChan, stopChan, resultChan := hex.CreatePatChecker()
	defer func() { stopChan <- struct{}{} }()

	_, a, err := alphaBeta(3, state, -abInit, abInit, gridChan, resultChan)

	if err != nil {
		log.Println(err)
	}

	if a == nil {
		// "Random" - TODO
		possibleActions := state.GetPossibleActions()
		return possibleActions[0].(*hex.Action)
	}

	return a
}

func alphaBeta(depth int, state *hex.State, alpha, beta float64, gridChan chan []uint64, resultChan chan [2][]int) (float64, *hex.Action, error) {
	if goal, _ := state.IsGoalState(false); goal {
		// Tha game has ended - the player who's turn it is has lost
		return -won, nil, nil
	}
	if depth <= 0 {
		val, err := eval(state, gridChan, resultChan)
		if err != nil {
			return 0, nil, err
		}
		return val, nil, nil
	}

	bestValue := -maxValue
	var bestState *hex.State

	possibleActions := state.GetPossibleActions()
	for _, a := range possibleActions {
		successor := state.GetSuccessorState(a).(hex.State)
		value, _, err := alphaBeta(depth-1, &successor, -beta, -alpha, gridChan, resultChan)
		if err != nil {
			return 0, nil, err
		}
		value = -value

		if value > bestValue {
			bestValue = value
			bestState = &successor
		}

		if bestValue > alpha {
			alpha = bestValue
			if alpha >= beta {
				// Prune
				break
			}
		}
	}

	var retAction *hex.Action
	if bestState != nil {
		retAction = state.GetTransitionAction(*bestState).(*hex.Action)
	}
	return bestValue, retAction, nil
}

func eval(state *hex.State, gridChan chan []uint64, resultChan chan [2][]int) (float64, error) {
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
	val := sample.getEstimatedValue()
	switch c := state.GetLastPlayer().Opponent(); c {
	case hex.Red:
		return val, nil
	case hex.Blue:
		return -val, nil
	default:
		return 0, fmt.Errorf("Invalid color %v", c)
	}
}
