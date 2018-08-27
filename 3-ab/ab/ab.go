package ab

import (
	"context"
	"fmt"
	"time"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
	"github.com/RdecKa/bachleor-thesis/common/tree"
)

const maxValue = 10000.0
const abInit = 1000.0
const won = 500.0

// AlphaBeta runs search with AB pruning to select the next action to be taken
func AlphaBeta(state *hex.State, timeToRun time.Duration, patFileName string) (*hex.Action, *tree.Tree) {
	gridChan, stopChan, resultChan := hex.CreatePatChecker(patFileName)
	defer func() { stopChan <- struct{}{} }()

	var selectedAction, a *hex.Action
	var rootNode *tree.Node
	var err error

	timeout := timeToRun
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	boardSize := state.GetSize()
	for depth := 2; depth < boardSize*boardSize; depth += 2 {
		fmt.Printf("Starting AB on depth %d\n", depth)

		transpositionTable := make(map[string]float64)
		_, a, rootNode, err = alphaBeta(ctx, depth, state, nil, -abInit, abInit, gridChan, resultChan, transpositionTable)

		if err != nil {
			fmt.Println(err)
			break
		}

		// If an action was not found in a shallower search, it will not
		// be found in a deeper search
		if a == nil {
			break
		}

		selectedAction = a
		fmt.Printf("Selected action: %v\n", selectedAction)
	}

	// Cancel the Context
	cancel()

	return selectedAction, tree.NewTree(rootNode)
}

func alphaBeta(ctx context.Context, depth int, state *hex.State, lastAction *hex.Action,
	alpha, beta float64, gridChan chan []uint32, resultChan chan [2][]int,
	transpositionTable map[string]float64) (float64, *hex.Action, *tree.Node, error) {

	// End recursion on timeout
	select {
	case <-ctx.Done():
		return 0, nil, nil, ctx.Err()
	default:
	}

	if val, ok := transpositionTable[state.GetMapKey()]; ok {
		// Current state was already investigated
		leaf := tree.NewNode(CreateAbNodeValue(state, val, "transT"))
		return val, lastAction, leaf, nil
	}
	if goal, _ := state.IsGoalState(false); goal {
		// The game has ended - the player who's turn it is has lost
		transpositionTable[state.GetMapKey()] = -won
		leaf := tree.NewNode(CreateAbNodeValue(state, -won, "goal"))
		return -won, lastAction, leaf, nil
	}
	if depth <= 0 {
		val, err := eval(state, gridChan, resultChan)
		if err != nil {
			return 0, nil, nil, err
		}
		transpositionTable[state.GetMapKey()] = val
		leaf := tree.NewNode(CreateAbNodeValue(state, val, "depth"))
		return val, lastAction, leaf, nil
	}

	bestValue := -maxValue
	var bestState *hex.State

	possibleActions := state.GetPossibleActions()
	nodeChildren := make([]*tree.Node, 0, len(possibleActions))
	comment := ""
	for _, a := range possibleActions {
		// End recursion on timeout
		select {
		case <-ctx.Done():
			return 0, nil, nil, ctx.Err()
		default:
		}

		successor := state.GetSuccessorState(a).(hex.State)
		value, _, childNode, err := alphaBeta(ctx, depth-1, &successor, a.(*hex.Action), -beta, -alpha, gridChan, resultChan, transpositionTable)
		if err != nil {
			return 0, nil, nil, err
		}
		value = -value

		if value > bestValue {
			bestValue = value
			bestState = &successor
		}

		nodeChildren = append(nodeChildren, childNode)

		if bestValue > alpha {
			alpha = bestValue
			if alpha >= beta {
				// Prune
				comment = "P"
				break
			}
		}
	}

	var retAction *hex.Action
	if bestState != nil {
		retAction = state.GetTransitionAction(*bestState).(*hex.Action)
	}
	transpositionTable[state.GetMapKey()] = bestValue
	node := tree.NewNode(CreateAbNodeValue(state, bestValue, comment))
	node.SetChildren(nodeChildren)
	return bestValue, retAction, node, nil
}

func eval(state *hex.State, gridChan chan []uint32, resultChan chan [2][]int) (float64, error) {
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
