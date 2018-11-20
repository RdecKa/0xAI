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

	gridChan, patChan, stopChan, resultChan := CreatePatChecker("patterns.txt")
	defer func() { stopChan <- struct{}{} }()

	sam := strings.Trim(state.GenSample(q, gridChan, patChan, resultChan), "\n")

	expected := expectedRed + "\n" + expectedBlue

	if sam != expected {
		t.Fatalf("Got\n'%s',\nexpected\n'%s'", sam, expected)
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
	redP := "1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	blueP := "1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	expectedRed := "-0.500000,2,1,1,1,11,10,1,1,1,1," + redP + "," + blueP
	expectedBlue := "0.500000,2,0,1,1,10,11,1,1,1,1," + blueP + "," + redP

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
	rr := "4,0,0,0,1,0,0,0,0,0,0,2,0,0,0,0,0,0,0,0,0,0,0,0" // Red attributes
	rt := "4,0,0,0,1,0,0,0,0,0,0,0,2,0,0,0,0,0,0,0,0,0,0,0" // Red attributes transposed
	bb := "4,2,1,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0" // Blue attributes
	bt := "4,1,1,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0" // Blue attributes transposed
	expectedRed := "-0.500000,8,1,7,5,17,23,4,3,5,3," + rr + "," + bb
	expectedBlue := "0.500000,8,0,5,7,23,17,3,5,3,4," + bt + "," + rt

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
	rr := "5,0,0,0,0,0,0,0,0,0,0,1,1,2,0,0,0,0,2,1,0,0,0,0"
	rt := "5,0,0,0,0,0,0,0,0,0,0,1,1,2,0,0,0,0,1,2,0,0,0,0"
	bb := "4,0,1,0,0,0,0,1,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0"
	bt := "4,0,1,0,0,1,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	expectedRed := "0.500000,9,0,4,5,16,20,4,4,3,4," + rr + "," + bb
	expectedBlue := "-0.500000,9,1,5,4,20,16,4,3,4,4," + bt + "," + rt

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
	rr := "5,0,0,0,0,0,0,2,0,0,0,0,1,1,0,0,0,0,0,0,0,0,0,0"
	rt := "5,0,0,0,0,2,0,0,0,0,0,1,0,1,0,0,0,0,0,0,0,0,0,0"
	bb := "5,1,1,0,0,0,0,2,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0"
	bt := "5,0,1,1,0,2,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	expectedRed := "-0.500000,10,1,5,8,20,19,5,4,4,5," + rr + "," + bb
	expectedBlue := "0.500000,10,0,8,5,19,20,5,4,4,5," + bt + "," + rt

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
	rr := "3,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	rt := "3,0,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	bb := "2,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	bt := "2,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	expectedRed := "0.500000,5,0,4,3,12,8,4,3,2,3," + rr + "," + bb
	expectedBlue := "-0.500000,5,1,3,4,8,12,3,2,3,4," + bt + "," + rt

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
	rr := "3,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	rt := "3,0,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	bb := "3,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	bt := "3,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	expectedRed := "-0.500000,6,1,4,7,12,11,4,3,3,4," + rr + "," + bb
	expectedBlue := "0.500000,6,0,7,4,11,12,4,3,3,4," + bt + "," + rt

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
	rr := "4,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	rt := "4,0,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	bb := "3,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	bt := "3,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
	expectedRed := "0.500000,7,0,23,17,21,20,5,4,3,4," + rr + "," + bb
	expectedBlue := "-0.500000,7,1,17,23,20,21,4,3,4,5," + bt + "," + rt

	test(t, 11, actions, expectedRed, expectedBlue)
}

/*
. . . . . . . . . . .
 . . . r . . . r r . .
  . . r r . . r r . b .
   . . . . . . . b b b b
    . . . r . r r b b b b
     . . r r . r r b b b b
      . . r . . . . . b b .
       . . . . . . . . . . .
        . . . . . . . . . . .
         . . . . . . . . . . .
          . . . . . . . . . . .
*/
func TestGenSample8(t *testing.T) {
	actions := []*Action{
		NewAction(3, 1, Red),
		NewAction(9, 2, Blue),
		NewAction(2, 2, Red),
		NewAction(7, 3, Blue),
		NewAction(3, 2, Red),
		NewAction(8, 3, Blue),
		NewAction(7, 1, Red),
		NewAction(9, 3, Blue),
		NewAction(8, 1, Red),
		NewAction(10, 3, Blue),
		NewAction(6, 2, Red),
		NewAction(7, 4, Blue),
		NewAction(7, 2, Red),
		NewAction(8, 4, Blue),
		NewAction(3, 4, Red),
		NewAction(9, 4, Blue),
		NewAction(2, 5, Red),
		NewAction(10, 4, Blue),
		NewAction(3, 5, Red),
		NewAction(7, 5, Blue),
		NewAction(2, 6, Red),
		NewAction(8, 5, Blue),
		NewAction(5, 4, Red),
		NewAction(9, 5, Blue),
		NewAction(6, 4, Red),
		NewAction(10, 5, Blue),
		NewAction(5, 5, Red),
		NewAction(8, 6, Blue),
		NewAction(6, 5, Red),
		NewAction(9, 6, Blue),
	}
	rr := "15,1,0,1,0,1,0,0,0,0,0,6,6,6,0,0,0,2,2,2,7,1,1,1"
	rt := "15,1,0,1,0,0,0,1,0,0,0,6,6,6,0,0,0,2,2,2,7,1,1,1"
	bb := "15,0,0,0,0,0,0,0,0,0,0,10,11,9,6,7,4,16,11,12,16,7,6,5"
	bt := "15,0,0,0,0,0,0,0,0,0,0,11,10,9,7,6,4,16,12,11,16,7,5,6"
	expectedRed := "-0.500000,30,1,45,55,44,16,6,7,5,4," + rr + "," + bb
	expectedBlue := "0.500000,30,0,55,45,16,44,4,5,7,6," + bt + "," + rt

	test(t, 11, actions, expectedRed, expectedBlue)
}
