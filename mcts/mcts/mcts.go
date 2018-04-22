// Package mcts provides Monte Carlo Tree Search
package mcts

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/RdecKa/mcts/game"
	"github.com/RdecKa/mcts/tree"
)

// -------------------------
// |     mctsNodeValue     |
// -------------------------

type mctsNodeValue struct {
	state game.State // state that this node represents
	n     uint       // how many times this node was visited
	q     float64    // estimated value of State state
}

func (mnv mctsNodeValue) String() string {
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
}

func (mcts *MCTS) String() string {
	return mcts.mcTree.String()
}

// InitMCTS initializes MCTS (State s is inserted in the root)
func InitMCTS(s game.State, c float64) *MCTS {
	node := createMCTSNode(s)
	mctsTree := tree.NewTree(node)
	rand.Seed(time.Now().UTC().UnixNano())
	return &MCTS{mctsTree, c}
}

// createMCTSNode creates new node with value {state=s, n=0, q=0}
func createMCTSNode(s game.State) *tree.Node {
	value := mctsNodeValue{s, 0, 0}
	node := tree.NewNode(&value)
	return node
}

// RunIteration runs one iteration of MCTS
func (mcts *MCTS) RunIteration() {
	mcts.selExpPlayBack(mcts.mcTree.GetRoot())
}

// selExpPlayBack performs one iteration of MCTS
// Phases of MCTS:
// 	selection: recursively call itself on the node's child with the highest
//		UCT value
// 	expansion: expand leaf node that was reached by recursive call
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
		// 	Old leaf node
		nodeValue.updateNodeValues(score)
		// 	New leaf node (if it was added in expansion phase)
		if newNode != nil {
			newNode.GetValue().(*mctsNodeValue).updateNodeValues(score)
		}

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
	score := mcts.selExpPlayBack(bestNode)

	// Update N and Q values (backpropagation)
	nodeValue.updateNodeValues(score)

	return score
}

// expansion finds all possible successor states and adds them as child nodes
// of Node node
func (mcts *MCTS) expansion(node *tree.Node) {
	nodeValue := node.GetValue().(*mctsNodeValue)
	state := nodeValue.state
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
	possibleActions := state.GetPossibleActions()
	if possibleActions == nil || len(possibleActions) <= 0 {
		return state.EvaluateGoalState()
	}
	randomAction := possibleActions[rand.Intn(len(possibleActions))]
	return playoutFromState(state.GetSuccessorState(randomAction))
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
