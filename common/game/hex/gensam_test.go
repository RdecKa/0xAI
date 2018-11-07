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

/*
. . . . . .
 . . . . b .
  . r . . . .
   . . . . . .
    . . . . . .
     . . . . . .
*/
func TestGenSample1(t *testing.T) {
	actions := []*Action{
		NewAction(1, 2, Red),
		NewAction(4, 1, Blue),
	}
	redP := "1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	blueP := "1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	expectedRed := "-0.500000,2,1,1,1,1,1,1," + redP + "," + blueP
	expectedBlue := "0.500000,2,0,1,1,1,1,1," + blueP + "," + redP

	test(t, 6, actions, expectedRed, expectedBlue)
}

/*
. . . . . r r .
 . . . b . . . .
  . . . . . . . .
   . . b . r r . .
    . . . b . . . .
     . b . . . . . .
      . . . . . . . .
       . . . . . . . .
*/
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
	rr := "4,0,0,0,1,0,0,0,0,0,0,2,0,0,0,0,0,0,0,0" // Red attributes
	rt := "4,0,0,0,1,0,0,0,0,0,0,0,2,0,0,0,0,0,0,0" // Red attributes transposed
	bb := "4,2,1,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0" // Blue attributes
	bt := "4,1,1,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0" // Blue attributes transposed
	expectedRed := "-0.500000,8,1,2,4,3,5,3," + rr + "," + bb
	expectedBlue := "0.500000,8,0,2,3,5,3,4," + bt + "," + rt

	test(t, 8, actions, expectedRed, expectedBlue)
}

/*
. . . . . . . .
 . . . . . . r .
  . . . . r r b .
   . b . r b . . .
    . . b r . . . .
     . . . . . . . .
      . . . . . . . .
       . . . . . . . .
*/
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
	rr := "5,0,0,0,0,0,0,0,0,0,0,1,1,2,0,0,0,0,2,1"
	rt := "5,0,0,0,0,0,0,0,0,0,0,1,1,2,0,0,0,0,1,2"
	bb := "4,0,1,0,0,0,0,1,0,0,1,0,0,0,0,0,0,0,0,0"
	bt := "4,0,1,0,0,1,0,0,1,0,0,0,0,0,0,0,0,0,0,0"
	expectedRed := "0.500000,9,0,0,4,4,3,4," + rr + "," + bb
	expectedBlue := "-0.500000,9,1,0,4,3,4,4," + bt + "," + rt

	test(t, 8, actions, expectedRed, expectedBlue)
}

/*
. . . . . . . .
 . . . . . . r .
  . . . . . r b .
   . b . r b . . .
    . . b r . . . .
     b r . . . . . .
      . . . . . . . .
       . . . . . . . .
*/
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
	rr := "5,0,0,0,0,0,0,2,0,0,0,0,1,1,0,0,0,0,0,0"
	rt := "5,0,0,0,0,2,0,0,0,0,0,1,0,1,0,0,0,0,0,0"
	bb := "5,1,1,0,0,0,0,2,0,0,1,0,0,0,0,0,0,0,0,0"
	bt := "5,0,1,1,0,2,0,0,1,0,0,0,0,0,0,0,0,0,0,0"
	expectedRed := "-0.500000,10,1,3,5,4,4,5," + rr + "," + bb
	expectedBlue := "0.500000,10,0,3,5,4,4,5," + bt + "," + rt

	test(t, 8, actions, expectedRed, expectedBlue)
}

/*
. . . r .
 . . b . .
  b . r . .
   . . . . .
    . r . . .
*/
func TestGenSample5(t *testing.T) {
	actions := []*Action{
		NewAction(3, 0, Red),
		NewAction(0, 2, Blue),
		NewAction(1, 4, Red),
		NewAction(2, 1, Blue),
		NewAction(2, 2, Red),
	}
	rr := "3,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	rt := "3,0,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0"
	bb := "2,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	bt := "2,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	expectedRed := "0.500000,5,0,0,4,3,2,3," + rr + "," + bb
	expectedBlue := "-0.500000,5,1,0,3,2,3,4," + bt + "," + rt

	test(t, 5, actions, expectedRed, expectedBlue)
}

/*
. . . r .
 . . b . .
  b . r . .
   . . . . .
    . r . . b
*/
func TestGenSample6(t *testing.T) {
	actions := []*Action{
		NewAction(3, 0, Red),
		NewAction(0, 2, Blue),
		NewAction(1, 4, Red),
		NewAction(2, 1, Blue),
		NewAction(2, 2, Red),
		NewAction(4, 4, Blue),
	}
	rr := "3,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	rt := "3,0,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0"
	bb := "3,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	bt := "3,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	expectedRed := "-0.500000,6,1,4,4,3,3,4," + rr + "," + bb
	expectedBlue := "0.500000,6,0,4,4,3,3,4," + bt + "," + rt

	test(t, 5, actions, expectedRed, expectedBlue)
}

/*
. . . r . . . . . . .
 . . b . . . . . . . .
  b . r . . . . . . . .
   . . . . . . . . . . .
    . r . . b . . . . . .
     . . . . . . . . . . .
      . . . . . . . . . . .
       . . . . . . . . . . .
        . . . . . . . . . . .
         . . . . . . . . . . .
          r . . . . . . . . . .
*/
func TestGenSample7(t *testing.T) {
	actions := []*Action{
		NewAction(3, 0, Red),
		NewAction(0, 2, Blue),
		NewAction(1, 4, Red),
		NewAction(2, 1, Blue),
		NewAction(2, 2, Red),
		NewAction(4, 4, Blue),
		NewAction(0, 10, Red),
	}
	rr := "4,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	rt := "4,0,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0"
	bb := "3,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	bt := "3,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	expectedRed := "0.500000,7,0,5,5,4,3,4," + rr + "," + bb
	expectedBlue := "-0.500000,7,1,5,4,3,4,5," + bt + "," + rt

	test(t, 11, actions, expectedRed, expectedBlue)
}
