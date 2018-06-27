package hex

import (
	"fmt"
)

// GenSample returns a string representation of a single learning sample (output, attributes...)
func (s State) GenSample(q float64, gridChan chan []uint64, resultChan chan [2][]int) string {
	gridChan <- s.GetCopyGrid()

	if s.lastPlayer == Blue {
		// Always store the Q value for the red player
		q = -q
	}
	red, blue, _ := s.GetNumOfStones()

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

	// Last two numbers in each c are numbers of occupied rows and columns for each player
	patCount := <-resultChan
	var patCountS string
	for _, p := range patCount {
		for _, c := range p {
			patCountS += fmt.Sprintf("%d,", c)
		}
	}

	return fmt.Sprintf("%f,%d,%d,%s\n", q, result, red+blue, patCountS)
}
