package ab

import (
	"fmt"
	"log"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

const bigNumber = 10000.0

// AlphaBeta runs search with AB pruning to select the next action to be taken
func AlphaBeta(state *hex.State) *hex.Action {
	gridChan, stopChan, resultChan := hex.CreatePatChecker()
	defer func() { stopChan <- struct{}{} }()

	val, a, err := alphaBeta(4, state, -bigNumber, bigNumber, gridChan, resultChan)

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Selected action %s with value %f.\n", a, val)

	if a == nil {
		// "Random" - TODO
		possibleActions := state.GetPossibleActions()
		return possibleActions[0].(*hex.Action)
	}

	return a
}

func alphaBeta(depth int, state *hex.State, alpha, beta float64, gridChan chan []uint64, resultChan chan [2][]int) (float64, *hex.Action, error) {
	goal, _ := state.IsGoalState(false)
	if goal || depth <= 0 {
		val, err := eval(state, gridChan, resultChan)
		if err != nil {
			return 0, nil, err
		}
		fmt.Printf("Investigated (depth %d), value %f:\n", depth, val)
		fmt.Printf("%s", state)
		return val, nil, nil
	}
	bestValue := -bigNumber
	var bestState hex.State

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
			bestState = successor
		}

		if bestValue > alpha {
			alpha = bestValue
			if alpha >= beta {
				fmt.Printf("Prune because %f >= %f\n", alpha, beta)
				break
			}
		}
	}

	fmt.Printf("Investigated (depth %d), value %f:\n", depth, bestValue)
	fmt.Printf("%s", state)

	return bestValue, state.GetTransitionAction(bestState).(*hex.Action), nil
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
