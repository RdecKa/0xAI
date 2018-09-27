package ab

import (
	"fmt"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

// ---------------------
// |     NodeValue     |
// ---------------------

// NodeValue stores a hex state and its estimated value, obtained by negamax
// algorithm with AB pruning. comment is used for debugging purposes
type NodeValue struct {
	state   *hex.State
	value   float64
	comment string
}

func (anv NodeValue) String() string {
	s := anv.state.String()
	s += fmt.Sprintf("(%f)\n", anv.value)
	return s
}

// CreateAbNodeValue creates a new NodeValue with given state and its estimated
// value
func CreateAbNodeValue(state *hex.State, value float64, comment string) *NodeValue {
	return &NodeValue{state, value, comment}
}

// -------------------------
// |     RootNodeValue     |
// -------------------------

// RootNodeValue stores a pointer to the actual AB tree together with some
// additional information about the game
/*type RootNodeValue struct {
	root  *tree.Node
	state *hex.State
	size  int
}

// CreateAbRootNodeValue creates a root node of the AB search tree
func CreateAbRootNodeValue(root *tree.Node, state *hex.State, size int) *RootNodeValue {
	return &RootNodeValue{root, state, size}
}*/
