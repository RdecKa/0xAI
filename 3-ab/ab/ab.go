package ab

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
	"github.com/RdecKa/bachleor-thesis/common/tree"
)

var maxValue = math.Inf(0)
var abInit = maxValue
var won = abInit

// AlphaBeta runs search with AB pruning to select the next action to be taken.
// In addition to the selected action it returns the tree that was constructed
// during the last AB search (if wanted).
func AlphaBeta(state *hex.State, timeToRun time.Duration, createTree bool,
	gridChan chan []uint32, resultChan chan [2][]int,
	getEstimatedValue func(s *Sample) float64) (*hex.Action, *tree.Tree) {

	var val float64
	var selectedAction, a *hex.Action
	var rootNode, rn *tree.Node
	var err error
	var oldTransitionTable map[string]float64

	timeout := timeToRun
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	boardSize := state.GetSize()
	for depthLimit := 2; depthLimit < boardSize*boardSize; depthLimit += 2 {
		// fmt.Printf("Starting AB on depth %d\n", depthLimit)

		transpositionTable := make(map[string]float64)
		val, a, rn, err = alphaBeta(ctx, 0, depthLimit, state, nil, -abInit, abInit, gridChan, resultChan, transpositionTable, oldTransitionTable, createTree, getEstimatedValue)
		oldTransitionTable = transpositionTable

		if err != nil {
			// fmt.Println(err)
			break
		}

		// If an action was not found in a shallower search, it will not
		// be found in a deeper search
		if a == nil {
			break
		}

		selectedAction = a
		rootNode = rn
		// fmt.Printf("Selected action: %v\n", selectedAction)

		// If the game is decided there is no need to continue with the search
		if math.IsInf(val, 0) {
			break
		}
	}

	// Cancel the Context
	cancel()

	// Create a tree for debuginng purposes
	var searchTree *tree.Tree
	if createTree && rootNode != nil {
		searchTree = tree.NewTree(rootNode)
	}

	return selectedAction, searchTree
}

func alphaBeta(ctx context.Context, depth, depthLimit int, state *hex.State,
	lastAction *hex.Action, alpha, beta float64, gridChan chan []uint32, resultChan chan [2][]int,
	transpositionTable, oldTransitionTable map[string]float64,
	createTree bool, getEstimatedValue func(s *Sample) float64) (float64, *hex.Action, *tree.Node, error) {

	// End recursion on timeout
	select {
	case <-ctx.Done():
		return 0, nil, nil, ctx.Err()
	default:
	}

	var leaf *tree.Node

	if val, ok := transpositionTable[state.GetMapKey()]; ok {
		// Current state was already investigated
		if createTree {
			leaf = tree.NewNode(CreateAbNodeValue(state, val, "TT"))
		}
		return val, lastAction, leaf, nil
	}
	if goal, _ := state.IsGoalState(false); goal {
		// The game has ended - the player who's turn it is has lost
		transpositionTable[state.GetMapKey()] = -won
		if createTree {
			leaf = tree.NewNode(CreateAbNodeValue(state, -won, "G"))
		}
		return -won, lastAction, leaf, nil
	}
	if depth >= depthLimit {
		val, err := eval(state, gridChan, resultChan, getEstimatedValue)
		if err != nil {
			return 0, nil, nil, err
		}
		transpositionTable[state.GetMapKey()] = val
		if createTree {
			leaf = tree.NewNode(CreateAbNodeValue(state, val, "D"))
		}
		return val, lastAction, leaf, nil
	}

	bestValue := -maxValue
	var bestState *hex.State

	possibleActions := state.GetPossibleActions()
	possibleActions = orderMoves(state, possibleActions, oldTransitionTable, true)

	var nodeChildren []*tree.Node
	if createTree {
		nodeChildren = make([]*tree.Node, 0, len(possibleActions))
	}
	comment := ""
	for _, a := range possibleActions {
		// End recursion on timeout
		select {
		case <-ctx.Done():
			return 0, nil, nil, ctx.Err()
		default:
		}

		successor := state.GetSuccessorState(a).(hex.State)
		value, _, childNode, err := alphaBeta(ctx, depth+1, depthLimit,
			&successor, a.(*hex.Action), -beta, -alpha, gridChan, resultChan,
			transpositionTable, oldTransitionTable, createTree, getEstimatedValue)
		if err != nil {
			return 0, nil, nil, err
		}
		value = -value

		if value > bestValue || bestState == nil {
			bestValue = value
			bestState = &successor
		}

		if createTree {
			nodeChildren = append(nodeChildren, childNode)
		}

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
		retAction = bestState.GetLastAction()
	}
	transpositionTable[state.GetMapKey()] = bestValue

	var node *tree.Node
	if createTree {
		node = tree.NewNode(CreateAbNodeValue(state, bestValue, comment))
		node.SetChildren(nodeChildren)
	}

	return bestValue, retAction, node, nil
}

// eval returns the estimated value of a sample
func eval(state *hex.State, gridChan chan []uint32, resultChan chan [2][]int,
	getEstimatedValue func(s *Sample) float64) (float64, error) {
	gridChan <- state.GetCopyGrid()
	patCount := <-resultChan

	args := &[]interface{}{*state, patCount}
	sample := Sample{
		num_stones: hex.AttrNumStones.GetAttributeValue(args),
		lp:         hex.AttrLastPlayer.GetAttributeValue(args),

		sdtc_r: hex.AttrDistanceToCenterRed.GetAttributeValue(args),
		sdtc_b: hex.AttrDistanceToCenterBlue.GetAttributeValue(args),
		rec_r:  hex.AttrReachableRed.GetAttributeValue(args),
		rec_b:  hex.AttrReachableBlue.GetAttributeValue(args),

		occ_red_rows:  hex.AttrOccRedRows.GetAttributeValue(args),
		occ_red_cols:  hex.AttrOccRedCols.GetAttributeValue(args),
		occ_blue_rows: hex.AttrOccBlueRows.GetAttributeValue(args),
		occ_blue_cols: hex.AttrOccBlueCols.GetAttributeValue(args),

		red_p0:  hex.AttrPatCountRed0.GetAttributeValue(args),
		red_p1:  hex.AttrPatCountRed1.GetAttributeValue(args),
		red_p2:  hex.AttrPatCountRed2.GetAttributeValue(args),
		red_p3:  hex.AttrPatCountRed3.GetAttributeValue(args),
		red_p4:  hex.AttrPatCountRed4.GetAttributeValue(args),
		red_p5:  hex.AttrPatCountRed5.GetAttributeValue(args),
		red_p6:  hex.AttrPatCountRed6.GetAttributeValue(args),
		red_p7:  hex.AttrPatCountRed7.GetAttributeValue(args),
		red_p8:  hex.AttrPatCountRed8.GetAttributeValue(args),
		red_p9:  hex.AttrPatCountRed9.GetAttributeValue(args),
		red_p10: hex.AttrPatCountRed10.GetAttributeValue(args),
		red_p11: hex.AttrPatCountRed11.GetAttributeValue(args),
		red_p12: hex.AttrPatCountRed12.GetAttributeValue(args),
		red_p13: hex.AttrPatCountRed13.GetAttributeValue(args),
		red_p14: hex.AttrPatCountRed14.GetAttributeValue(args),
		red_p15: hex.AttrPatCountRed15.GetAttributeValue(args),
		red_p16: hex.AttrPatCountRed16.GetAttributeValue(args),
		red_p17: hex.AttrPatCountRed17.GetAttributeValue(args),
		red_p18: hex.AttrPatCountRed18.GetAttributeValue(args),
		red_p19: hex.AttrPatCountRed19.GetAttributeValue(args),
		red_p20: hex.AttrPatCountRed20.GetAttributeValue(args),
		red_p21: hex.AttrPatCountRed21.GetAttributeValue(args),
		red_p22: hex.AttrPatCountRed22.GetAttributeValue(args),
		red_p23: hex.AttrPatCountRed23.GetAttributeValue(args),

		blue_p0:  hex.AttrPatCountBlue0.GetAttributeValue(args),
		blue_p1:  hex.AttrPatCountBlue1.GetAttributeValue(args),
		blue_p2:  hex.AttrPatCountBlue2.GetAttributeValue(args),
		blue_p3:  hex.AttrPatCountBlue3.GetAttributeValue(args),
		blue_p4:  hex.AttrPatCountBlue4.GetAttributeValue(args),
		blue_p5:  hex.AttrPatCountBlue5.GetAttributeValue(args),
		blue_p6:  hex.AttrPatCountBlue6.GetAttributeValue(args),
		blue_p7:  hex.AttrPatCountBlue7.GetAttributeValue(args),
		blue_p8:  hex.AttrPatCountBlue8.GetAttributeValue(args),
		blue_p9:  hex.AttrPatCountBlue9.GetAttributeValue(args),
		blue_p10: hex.AttrPatCountBlue10.GetAttributeValue(args),
		blue_p11: hex.AttrPatCountBlue11.GetAttributeValue(args),
		blue_p12: hex.AttrPatCountBlue12.GetAttributeValue(args),
		blue_p13: hex.AttrPatCountBlue13.GetAttributeValue(args),
		blue_p14: hex.AttrPatCountBlue14.GetAttributeValue(args),
		blue_p15: hex.AttrPatCountBlue15.GetAttributeValue(args),
		blue_p16: hex.AttrPatCountBlue16.GetAttributeValue(args),
		blue_p17: hex.AttrPatCountBlue17.GetAttributeValue(args),
		blue_p18: hex.AttrPatCountBlue18.GetAttributeValue(args),
		blue_p19: hex.AttrPatCountBlue19.GetAttributeValue(args),
		blue_p20: hex.AttrPatCountBlue20.GetAttributeValue(args),
		blue_p21: hex.AttrPatCountBlue21.GetAttributeValue(args),
		blue_p22: hex.AttrPatCountBlue22.GetAttributeValue(args),
		blue_p23: hex.AttrPatCountBlue23.GetAttributeValue(args),
	}
	val := getEstimatedValue(&sample)

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

// GetEstimateFunction returns a function that will be used for evaluating
// states
func GetEstimateFunction(subtype string) func(s *Sample) float64 {
	switch subtype {
	case "abDT":
		return getEstimatedValueDT
	case "abLR":
		return getEstimatedValueLR
	default:
		panic(fmt.Errorf("Invalid AB subtype: %s", subtype))
	}
}
