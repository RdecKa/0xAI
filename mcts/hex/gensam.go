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
	return fmt.Sprintf("%f,%d,%d,%d\n", q, red, blue, empty)
}
