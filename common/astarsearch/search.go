// Package astarsearch provides A* search that figures out whether a solution
// exists or not
package astarsearch

import (
	"container/heap"

	"github.com/RdecKa/common/pq"
	"github.com/RdecKa/common/tree"
)

// -----------------
// |     State     |
// -----------------

// State represents a state in A* Search
type State interface {
	IsGoalState() bool
	GetEstimateToReachGoal() int
	GetSuccessorStates() []State
	String() string
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
	frontier      pq.PriorityQueue   // List of nodes to be expanded
	visitedStates map[State]struct{} // A list of states that have already been added to the frontier
}

func makeAStarNode(state State, pathFromStart int) *tree.Node {
	nodeValue := aStarNodeValue{state, pathFromStart, pathFromStart + state.GetEstimateToReachGoal()}
	return tree.NewNode(&nodeValue)
}

// InitSearch initializes A* Search with a game state initialState where the
// solution is searched for
func InitSearch(initialState State) *AStarSearch {
	startNode := makeAStarNode(initialState, 0)
	newTree := tree.NewTree(startNode)

	newFrontier := pq.New(50)                        // Create a new frontier
	heap.Push(newFrontier, pq.NewItem(0, startNode)) // Init frontier with the initial state

	visitedStates := make(map[State]struct{})

	aStar := &AStarSearch{newTree, *newFrontier, visitedStates}
	return aStar
}

// Search tries to find a goal state. If a soltuion exists, it returns true,
// otherwise it returns false
func (aStar *AStarSearch) Search() bool {
	for len(aStar.frontier) > 0 {
		// Pop frontier
		currentNode := heap.Pop(&aStar.frontier).(*tree.Node)

		// Get state from node
		nodeValue := currentNode.GetValue().(*aStarNodeValue)
		currentState := nodeValue.state

		_, ok := aStar.visitedStates[currentState]
		if ok {
			// State was already expanded, discard ot
			continue
		}

		if currentState.IsGoalState() {
			// Solution found
			return true
		}

		// Add the current state to the list of visited states
		aStar.visitedStates[currentState] = struct{}{}

		// Get the cost of the path from start to the current node
		pathFromStart := nodeValue.pathFromStart

		// Loop through all successor states (find all possibilities to continue
		// the chain)
		successorStates := currentState.GetSuccessorStates()
		for _, sucState := range successorStates {
			_, ok := aStar.visitedStates[sucState]
			if ok {
				// Successor state was already expanded, discard it
				continue
			}

			// Add a new node to the priority queue
			sucNode := makeAStarNode(sucState, pathFromStart+1)
			priority := sucNode.GetValue().(*aStarNodeValue).heuristicValue
			heap.Push(&aStar.frontier, pq.NewItem(priority, sucNode))
		}

		/*fmt.Println("Priority queue:")
		for _, el := range aStar.frontier {
			fmt.Println(el.GetValue().(*tree.Node).GetValue())
			fmt.Println(el.GetValue().(*tree.Node).GetValue().(*aStarNodeValue).state)
		}
		fmt.Println("End of priority queue")*/
	}

	return false
}

// same compares two tree nodes and returns true if they contain the same state
func same(a, b interface{}) bool {
	if a == b {
		return true
	}

	av := a.(*tree.Node).GetValue().(*aStarNodeValue)
	bv := b.(*tree.Node).GetValue().(*aStarNodeValue)

	return (av.state.String() == bv.state.String())
}
