package ab

import (
	"fmt"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

// NodeValue stores a hex state and its estimated value, obtained by negamax
// algorithm with AB pruning
type NodeValue struct {
	state *hex.State
	value float64
}

func (anv NodeValue) String() string {
	s := anv.state.String()
	s += fmt.Sprintf("(%f)\n", anv.value)
	return s
}

// CreateAbNodeValue creates a new NodeValue with given state and its estimated
// value
func CreateAbNodeValue(state *hex.State, value float64) *NodeValue {
	return &NodeValue{state, value}
}
