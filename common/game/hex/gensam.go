package hex

import (
	"fmt"
)

// GenSample returns a string representation of a single learning sample (output, attributes...)
func (s State) GenSample(q float64, gridChan chan []uint32, resultChan chan [2][]int) string {
	gridChan <- s.GetCopyGrid()

	if s.lastPlayer == Blue {
		// Always store the Q value for the red player
		q = -q
	}

	o := fmt.Sprintf("%f", q)

	patCount := <-resultChan
	args := &[]interface{}{s, patCount}
	for _, attr := range GenSamAttributes {
		o += ","
		o += fmt.Sprintf("%v", attr.GetAttributeValue(args))
	}
	o += "\n"
	return o
}

// GetHeaderCSV returns a string consisting of attribute names.
func GetHeaderCSV() string {
	o := "value"
	for _, attr := range GenSamAttributes {
		o += ","
		o += attr.GetAttributeName()
	}
	o += "\n"
	return o
}
