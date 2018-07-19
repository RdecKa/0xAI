// Package mcts provides Monte Carlo Tree Search
package mcts

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/RdecKa/bachleor-thesis/common/game"
	"github.com/RdecKa/bachleor-thesis/common/tree"
)

// -------------------------
// |     mctsNodeValue     |
// -------------------------

type mctsNodeValue struct {
	state game.State // state that this node represents
	n     uint       // how many times this node was visited
	q     float64    // estimated value of State state
}

func (mnv *mctsNodeValue) String() string {
	return fmt.Sprintf("%s (N: %d, Q: %f)", mnv.state, mnv.n, mnv.q)
}

// updateNodeValues increments N and calculates new average for Q
func (mnv *mctsNodeValue) updateNodeValues(score float64) {
	mnv.n++
	mnv.q += (score - mnv.q) / float64(mnv.n)
}

// ----------------
// |     MCTS     |
// ----------------

// MCTS represens Monte Carlo Tree Search
type MCTS struct {
	mcTree *tree.Tree // Monte Carlo tree
	c      float64    // exploration parameter
	minN   uint       // minimal number of visits of a node before it can be expanded
}

func (mcts *MCTS) String() string {
	return mcts.mcTree.String()
}

// InitMCTS initializes MCTS (State s is inserted in the root)
func InitMCTS(s game.State, c float64, minN uint) *MCTS {
	node := createMCTSNode(s)
	mctsTree := tree.NewTree(node)
	rand.Seed(time.Now().UTC().UnixNano())
	return &MCTS{mctsTree, c, minN}
}

// ContinueMCTSFromNode continues MCTS from Node node
func (mcts *MCTS) ContinueMCTSFromNode(node *tree.Node) *MCTS {
	mctsTree := tree.NewTree(node)
	return &MCTS{mctsTree, mcts.c, mcts.minN}
}

// createMCTSNode creates new node with value {state=s, n=0, q=0}
func createMCTSNode(s game.State) *tree.Node {
	value := mctsNodeValue{s, 0, 0}
	node := tree.NewNode(&value)
	return node
}

// GetInitialNode returns the node in which the search has began
func (mcts *MCTS) GetInitialNode() *tree.Node {
	return mcts.mcTree.GetRoot()
}

// RunMCTS executes iterations of MCTS for timeToRun, given initialised MCTS
func RunMCTS(mc *MCTS, workerID int, timeToRun time.Duration, boardSize int, outputFile, logFile *os.File, gridChan chan []uint64, resultChan chan [2][]int) ([]*tree.Node, error) {
	timer := time.NewTimer(timeToRun)

	timeOut := false
	for iterCount := 0; !timeOut; iterCount++ {
		select {
		case <-timer.C:
			fmt.Println("TIME OUT")
			timeOut = true
		default:
			if iterCount > 0 && iterCount%10000 == 0 {
				logFile.WriteString(fmt.Sprintf("Worker %d finished iteration %d\n", workerID, iterCount))
			}
			mc.RunIteration()
			break
		}
	}

	// Write input-output pairs for supervised machine learning, generate
	// new nodes to continue MCTS
	expCand, err := mc.GenSamples(outputFile, 100, gridChan, resultChan)
	if err != nil {
		return nil, err
	}
	return expCand, nil
}

// RunIteration runs one iteration of MCTS
func (mcts *MCTS) RunIteration() {
	mcts.selExpPlayBack(mcts.mcTree.GetRoot())
}

// selExpPlayBack performs one iteration of MCTS
// Phases of MCTS:
// 	selection: recursively call itself on the node's child with the highest
//		UCT value
// 	expansion: expand leaf node that was reached by recursive call (only if the
//		node has been visited often enough)
// 	playout: randomly select moves until goal state is reached (no possible
//		actions)
// 	backpropagation: update values on nodes on selected branch in the tree
func (mcts *MCTS) selExpPlayBack(node *tree.Node) float64 {
	children := node.GetChildren()
	nodeValue := node.GetValue().(*mctsNodeValue)

	if len(children) == 0 {
		// Leaf node reached, selection phase finished

		// Expansion phase
		mcts.expansion(node)
		// Select one of the new children (if there are any) and run playout
		// from there
		newChildren := node.GetChildren()
		var newNode *tree.Node
		if len(newChildren) > 0 {
			newNode = newChildren[rand.Intn(len(newChildren))]
		}

		// Playout phase
		var score float64
		if newNode != nil {
			score = mcts.playout(newNode)
		} else {
			score = mcts.playout(node)
		}

		// Backpropagation begins - update two last nodes:
		// 	New leaf node (if it was added in expansion phase)
		if newNode != nil {
			newNode.GetValue().(*mctsNodeValue).updateNodeValues(score)
			score = -score // Negate the value because node is newNode's opponent!
		}
		// 	Old leaf node
		nodeValue.updateNodeValues(score)

		return score
	}

	// Iterate through all children, find the one with the highest UCT value
	maxUCTValue := mcts.getUCTValue(children[0], nodeValue.n)
	bestNode := children[0]

	for i, child := range children[1:] {
		UCTValue := mcts.getUCTValue(child, nodeValue.n)
		if UCTValue > maxUCTValue {
			maxUCTValue = UCTValue
			bestNode = children[i+1] // +1 because i starts at 0, but the array with children[1]
		}
	}

	// Recursive call (selection)
	score := -mcts.selExpPlayBack(bestNode)

	// Update N and Q values (backpropagation)
	nodeValue.updateNodeValues(score)

	return score
}

// expansion finds all possible successor states and adds them as child nodes
// of Node node
func (mcts *MCTS) expansion(node *tree.Node) {
	nodeValue := node.GetValue().(*mctsNodeValue)
	state := nodeValue.state

	if g, _ := state.IsGoalState(false); g {
		// Do not expand goal states
		return
	}

	if nodeValue.n < mcts.minN {
		// Do not expand a node that has not been visited at least minN times
		return
	}

	possibleActions := state.GetPossibleActions()
	successorNodes := make([]*tree.Node, len(possibleActions))

	for i, action := range possibleActions {
		successorNodes[i] = createMCTSNode(state.GetSuccessorState(action))
	}
	node.SetChildren(successorNodes)
}

// playout starts playout phase of MCTS from Node node
func (mcts *MCTS) playout(node *tree.Node) float64 {
	nodeValue := node.GetValue().(*mctsNodeValue)
	state := nodeValue.state
	return playoutFromState(state)
}

// playoutFromState recursively performs a random action from the list of
// possible actions. After reaching a goal state it returns its value
func playoutFromState(state game.State) float64 {
	if g, _ := state.IsGoalState(false); g {
		return state.EvaluateGoalState()
	}
	possibleActions := state.GetPossibleActions()
	if possibleActions == nil || len(possibleActions) <= 0 {
		panic(fmt.Sprintf("Not in a goal state yet, but no action possible. Something is wrong."))
	}
	randomAction := possibleActions[rand.Intn(len(possibleActions))]
	return -playoutFromState(state.GetSuccessorState(randomAction))
}

// getUCTValue calculates UCT value of a Node node.
// Argument n represents N value of parent node (how many times parent node was
// visited)
func (mcts *MCTS) getUCTValue(node *tree.Node, parentN uint) float64 {
	// Assert that nodeValue is of type *mctsNodeValue
	nodeValue := node.GetValue().(*mctsNodeValue)

	if nodeValue.n == 0 {
		return math.MaxFloat64
	}

	return float64(nodeValue.q) +
		mcts.c*math.Sqrt(math.Log(float64(parentN))/float64(nodeValue.n))
}

// ContinueMCTSFromChild goes through all grandchildren nodes of the root node
// in the MC tree and finds the one that contains the given state. It returns
// MCTS from that node.
func (mcts *MCTS) ContinueMCTSFromChild(state game.State) *MCTS {
	nodeChildren := mcts.mcTree.GetRoot().GetChildren()
	for _, n := range nodeChildren {
		grandChildren := n.GetChildren()
		for _, g := range grandChildren {
			s := g.GetValue().(*mctsNodeValue).state.(game.State)
			if s.Same(state) {
				return mcts.ContinueMCTSFromNode(g)
			}
		}
	}
	// Not possible to continue previously started search, start from scratch.
	return InitMCTS(state, mcts.c, mcts.minN)
}

// GetBestRootChildState returns agame.State of the direct descendant of the
// root node in the MC tree that has the highest UCT value.
// It returns nul if no such state exists.
func (mcts *MCTS) GetBestRootChildState() game.State {
	rootChildren := mcts.mcTree.GetRoot().GetChildren()
	if len(rootChildren) == 0 {
		return nil
	}
	bestNode := rootChildren[0]

	for _, c := range rootChildren {
		if c.GetValue().(*mctsNodeValue).q > bestNode.GetValue().(*mctsNodeValue).q {
			bestNode = c
		}
	}

	return bestNode.GetValue().(*mctsNodeValue).state
}
