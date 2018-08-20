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

	/*result := 0 // 0 if game not finished, 1 if red wins, -1 if blue wins
	if s.IsGoalState() {
		if s.lastPlayer == Blue {
			result = -1
		} else if s.lastPlayer == Red {
			result = 1
		} else {
			panic(fmt.Sprintf("Unknown color '%s'\n", s.lastPlayer))
		}
	}*/

	// Last two numbers in each c are numbers of occupied rows and columns for each player
	patCount := <-resultChan
	var patCountS, occCountS string
	for _, player := range patCount {
		for _, c := range player[:len(player)-2] {
			patCountS += fmt.Sprintf("%d,", c)
		}
		for _, c := range player[len(player)-2:] {
			occCountS += fmt.Sprintf("%d,", c)
		}
	}

	st := fmt.Sprintf("%f,%d,%s%s", q, red+blue, occCountS, patCountS)
	st = st[0 : len(st)-1]
	return fmt.Sprintf("%s\n", st)
}

// GetHeaderCSV returns a string consisting attribute names.
func GetHeaderCSV() string {
	s := "value,num_stones,occ_red_rows,occ_red_cols,occ_blue_rows,occ_blue_cols"
	for n := 0; n <= 2; n++ {
		s += fmt.Sprintf(",red_p%d", n)
	}
	for n := 0; n <= 2; n++ {
		s += fmt.Sprintf(",blue_p%d", n)
	}
	return fmt.Sprintf("%s\n", s)
}
