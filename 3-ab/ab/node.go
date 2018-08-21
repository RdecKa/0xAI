package ab

import (
	"fmt"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

type AbNodeValue struct {
	state *hex.State
	value float64
}

func (anv AbNodeValue) String() string {
	s := anv.state.String()
	s += fmt.Sprintf("(%f)\n", anv.value)
	return s
}

func CreateAbNodeValue(state *hex.State, value float64) *AbNodeValue {
	return &AbNodeValue{state, value}
}
