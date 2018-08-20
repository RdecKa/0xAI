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
func AlphaBeta(state *hex.State, patFileName string) *hex.Action {
	gridChan, stopChan, resultChan := hex.CreatePatChecker(patFileName)
	defer func() { stopChan <- struct{}{} }()

	transpositionTable := make(map[string]float64)

	_, a, err := alphaBeta(6, state, nil, -abInit, abInit, gridChan, resultChan, transpositionTable)

	if err != nil {
		log.Println(err)
	}

	if a == nil {
		// "Random" - TODO
		fmt.Println("Choosing 'randomly'")
		possibleActions := state.GetPossibleActions()
		return possibleActions[0].(*hex.Action)
	}

	return a
}

func alphaBeta(depth int, state *hex.State, lastAction *hex.Action,
	alpha, beta float64, gridChan chan []uint64, resultChan chan [2][]int,
	transpositionTable map[string]float64) (float64, *hex.Action, error) {

	if val, ok := transpositionTable[state.GetMapKey()]; ok {
		// Current state was already investigated
		return val, lastAction, nil
	}
	if goal, _ := state.IsGoalState(false); goal {
		// The game has ended - the player who's turn it is has lost
		transpositionTable[state.GetMapKey()] = -won
		return -won, lastAction, nil
	}
	if depth <= 0 {
		val, err := eval(state, gridChan, resultChan)
		if err != nil {
			return 0, nil, err
		}
		transpositionTable[state.GetMapKey()] = val
		return val, lastAction, nil
	}

	bestValue := -maxValue
	var bestState *hex.State

	possibleActions := state.GetPossibleActions()
	for _, a := range possibleActions {
		successor := state.GetSuccessorState(a).(hex.State)
		value, _, err := alphaBeta(depth-1, &successor, a.(*hex.Action), -beta, -alpha, gridChan, resultChan, transpositionTable)
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
	transpositionTable[state.GetMapKey()] = bestValue
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
