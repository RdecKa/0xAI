package hex

import (
	"fmt"
)

// GenSample returns a string representation of two learning samples in format:
// (output, attributes...)
// First learning sample is a representation of a given State s, the second is a
// representation of the same state but with reversed roles of red and blue
// player
func (s State) GenSample(q float64, gridChan chan []uint32, patChan chan []int, resultChan chan [2][]int) string {
	gridChan <- s.GetCopyGrid()
	patChan <- nil

	if s.lastAction.c == Blue {
		// Always store the Q value for the red player
		q = -q
	}

	// o1 <- State s
	// o2 <- inversed State s
	o1 := fmt.Sprintf("%f", q)
	o2 := fmt.Sprintf("%f", -q)

	patCount := <-resultChan
	args := &[]interface{}{s, patCount}
	for _, attrPair := range GenSamAttributes {
		aVal := attrPair[0].GetAttributeValue(args)
		o1 += fmt.Sprintf(",%v", aVal)

		if attrPair[1] != nil {
			aVal = attrPair[1].GetAttributeValue(args)
		}
		o2 += fmt.Sprintf(",%v", aVal)
	}
	return o1 + "\n" + o2 + "\n"
}

// GetHeaderCSV returns a string consisting of attribute names.
func GetHeaderCSV() string {
	o := "value"
	for _, attrPair := range GenSamAttributes {
		o += ","
		o += attrPair[0].GetAttributeName()
	}
	o += "\n"
	return o
}
