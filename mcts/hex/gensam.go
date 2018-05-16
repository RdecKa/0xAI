package hex

import (
	"fmt"
)

// GenSample returns a string representation of a single learning sample (output, attributes...)
func (s State) GenSample(q float64) string {
	if s.lastPlayer == Blue {
		// Always store the Q value for the red player
		q = -q
	}
	red, blue, empty := s.GetNumOfStones()

	result := 0 // 0 if game not finished, 1 if red wins, -1 if blue wins
	if s.IsGoalState() {
		if s.lastPlayer == Blue {
			result = -1
		} else if s.lastPlayer == Red {
			result = 1
		} else {
			panic(fmt.Sprintf("Unknown color '%s'\n", s.lastPlayer))
		}
	}

	return fmt.Sprintf("%f,%d,%d,%d,%d\n", q, result, red, blue, empty)
}
