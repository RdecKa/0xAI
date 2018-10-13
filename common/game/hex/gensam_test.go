package hex

import (
	"strings"
	"testing"
)

func test(t *testing.T, size byte, actions []*Action, expectedRed, expectedBlue string) {
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

	expected := expectedRed + "\n" + expectedBlue

	if sam != expected {
		t.Fatalf("Got '%s', expected '%s'", sam, expected)
	}
}
func TestGenSample1(t *testing.T) {
	actions := []*Action{
		NewAction(1, 2, Red),
		NewAction(4, 1, Blue),
	}
	expectedRed := "-0.500000,2,1,1,1,1,1,0,0,0,0,1,0,0,0,0,1,1"
	expectedBlue := "0.500000,2,1,1,1,1,1,0,0,0,0,1,0,0,0,0,0,1"

	test(t, 6, actions, expectedRed, expectedBlue)
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
	expectedRed := "-0.500000,8,4,3,5,3,4,0,1,0,0,4,4,0,0,0,1,2"
	expectedBlue := "0.500000,8,3,5,3,4,4,4,0,0,0,4,0,1,0,0,0,2"

	test(t, 8, actions, expectedRed, expectedBlue)
}

func TestGenSample3(t *testing.T) {
	actions := []*Action{
		NewAction(6, 1, Red),
		NewAction(1, 3, Blue),
		NewAction(5, 2, Red),
		NewAction(2, 4, Blue),
		NewAction(4, 2, Red),
		NewAction(4, 3, Blue),
		NewAction(3, 3, Red),
		NewAction(6, 2, Blue),
		NewAction(3, 4, Red),
	}
	expectedRed := "0.500000,9,4,4,3,4,5,0,0,0,0,4,1,0,1,1,0,0"
	expectedBlue := "-0.500000,9,4,3,4,4,4,1,0,1,1,5,0,0,0,0,1,0"

	test(t, 8, actions, expectedRed, expectedBlue)
}

func TestGenSample4(t *testing.T) {
	actions := []*Action{
		NewAction(6, 1, Red),
		NewAction(1, 3, Blue),
		NewAction(5, 2, Red),
		NewAction(2, 4, Blue),
		NewAction(1, 5, Red),
		NewAction(4, 3, Blue),
		NewAction(3, 3, Red),
		NewAction(6, 2, Blue),
		NewAction(3, 4, Red),
		NewAction(0, 5, Blue),
	}
	expectedRed := "-0.500000,10,5,4,4,5,5,0,0,2,0,5,2,0,2,1,1,3"
	expectedBlue := "0.500000,10,5,4,4,5,5,2,0,2,1,5,0,0,2,0,0,3"

	test(t, 8, actions, expectedRed, expectedBlue)
}
