// Package astarsearch provides A* search that figures out whether a solution
// exists or not
package astarsearch

import (
	"github.com/RdecKa/mcts/tree"
)

// -----------------
// |     State     |
// -----------------

// State represents a state in A* Search
type State interface {
	IsGoalState() bool
	GetEstimateToReachGoal() int
	GetSuccessorStates() []State
}

// --------------------------
// |     aStarNodeValue     |
// --------------------------

type aStarNodeValue struct {
	state          State // State in a game
	pathFromStart  int   // Path cost from initial state
	heuristicValue int   // pathFromStart plus estimated cost to a goal state
}

// -----------------------
// |     AStarSearch     |
// -----------------------

// AStarSearch represents a search tree for A* Search and its frontier
type AStarSearch struct {
	tree          *tree.Tree         // Search tree
	frontier      []*tree.Node       // List of nodes to be expanded
	visitedStates map[State]struct{} // A list of states that have already been added to the frontier
}

func makeAStarNode(state State, pathFromStart int) *tree.Node {
	nodeValue := aStarNodeValue{state, pathFromStart, state.GetEstimateToReachGoal()}
	return tree.NewNode(&nodeValue)
}

// InitSearch initializes A* Search with a game state initialState where the
// solution is searched for
func InitSearch(initialState State) *AStarSearch {
	startNode := makeAStarNode(initialState, 0)
	newTree := tree.NewTree(startNode)
	newFrontier := make([]*tree.Node, 1, 50) // Create a new frontier (initial size = 1, reserved size = 50)
	newFrontier[0] = startNode               // Init frontier with the initial state
	visitedStates := make(map[State]struct{})
	visitedStates[initialState] = struct{}{} // Mark initial state as visited
	aStar := &AStarSearch{newTree, newFrontier, visitedStates}
	return aStar
}

// Search tries to find a goal state. If a soltuion exists, it returns true,
// otherwise it returns false
func (aStar *AStarSearch) Search() bool {
	for len(aStar.frontier) > 0 {
		// Pop frontier
		currentNode := aStar.frontier[0]
		aStar.frontier = aStar.frontier[1:]

		// Get state from node
		nodeValue := currentNode.GetValue().(*aStarNodeValue)
		currentState := nodeValue.state
		if currentState.IsGoalState() {
			// Solution found
			return true
		}

		// Get the cost of the path from start to the current node
		pathFromStart := nodeValue.pathFromStart

		// Loop through all successor states (find all possibilities to continue
		// the chain)
		successorStates := currentState.GetSuccessorStates()
		for _, sucState := range successorStates {
			_, ok := aStar.visitedStates[sucState]
			if ok {
				// State sucState revisited, discard it
				continue
			}

			// Add a state to the list of visited states
			aStar.visitedStates[sucState] = struct{}{}

			sucNode := makeAStarNode(sucState, pathFromStart+1)

			// TODO: check if discarding revisited states implemented correctly
			// TODO: insert in an ordered list - now it works as BFS
			aStar.frontier = append(aStar.frontier, sucNode)
		}
	}

	return false
}