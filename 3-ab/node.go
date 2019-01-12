package ab

import (
	"fmt"

	"github.com/RdecKa/0xAI/common/game/hex"
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
