package hex

import (
	"strings"
	"testing"
)

func test(t *testing.T, size byte, actions []*Action, expected string) {
	state := NewState(size, Red)
	for _, a := range actions {
		s := state.GetSuccessorState(a).(State)
		state = &s
	}

	q := 0.5
	gridChan := make(chan []uint32, 1)
	resultChan := make(chan [2][]int, 1)
	stopChan := make(chan struct{}, 1)

	go patChecker("patterns.txt", gridChan, stopChan, resultChan)

	sam := strings.Trim(state.GenSample(q, gridChan, resultChan), "\n")

	if sam != expected {
		t.Fatalf("Got '%s', expected '%s'", sam, expected)
	}
}
func TestGenSample1(t *testing.T) {
	actions := []*Action{
		NewAction(1, 2, Red),
		NewAction(4, 1, Blue),
	}
	expected := "-0.500000,2,1,1,1,1,1,0,0,1,0,0"

	test(t, 6, actions, expected)
}

func TestGenSample2(t *testing.T) {
	actions := []*Action{
		NewAction(5, 0, Red),
		NewAction(3, 1, Blue),
		NewAction(6, 0, Red),
		NewAction(2, 3, Blue),
		NewAction(4, 3, Red),
		NewAction(3, 4, Blue),
		NewAction(5, 3, Red),
		NewAction(1, 5, Blue),
	}
	expected := "-0.500000,8,4,3,5,3,4,0,1,4,4,0"

	test(t, 8, actions, expected)
}
