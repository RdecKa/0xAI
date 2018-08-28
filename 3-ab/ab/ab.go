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
	var rootNode, rn *tree.Node
	var err error

	timeout := timeToRun
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	boardSize := state.GetSize()
	for depth := 2; depth < boardSize*boardSize; depth += 2 {
		fmt.Printf("Starting AB on depth %d\n", depth)

		transpositionTable := make(map[string]float64)
		_, a, rn, err = alphaBeta(ctx, depth, state, nil, -abInit, abInit, gridChan, resultChan, transpositionTable)

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
		rootNode = rn
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
	patCount := <-resultChan

	args := &[]interface{}{*state, patCount}
	sample := Sample{
		num_stones:    hex.AttrNumStones.GetAttributeValue(args),
		occ_red_rows:  hex.AttrOccRedRows.GetAttributeValue(args),
		occ_red_cols:  hex.AttrOccRedCols.GetAttributeValue(args),
		occ_blue_rows: hex.AttrOccBlueRows.GetAttributeValue(args),
		occ_blue_cols: hex.AttrOccBlueCols.GetAttributeValue(args),
		red_p0:        hex.AttrPatCountRed0.GetAttributeValue(args),
		red_p1:        hex.AttrPatCountRed1.GetAttributeValue(args),
		red_p2:        hex.AttrPatCountRed2.GetAttributeValue(args),
		blue_p0:       hex.AttrPatCountBlue0.GetAttributeValue(args),
		blue_p1:       hex.AttrPatCountBlue1.GetAttributeValue(args),
		blue_p2:       hex.AttrPatCountBlue2.GetAttributeValue(args),
	}
	val := sample.getEstimatedValue()

	// val is given from Red player's prospective
	switch c := state.GetLastPlayer().Opponent(); c {
	case hex.Red:
		return val, nil
	case hex.Blue:
		return -val, nil
	default:
		return 0, fmt.Errorf("Invalid color %v", c)
	}
}
