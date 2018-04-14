package main

import (
	"fmt"
	"strconv"

	"github.com/RdecKa/mcts/mcts"
)

// Dummy state
type dummyState struct {
	num   int
	depth int
}

func (s dummyState) String() string {
	return strconv.Itoa(s.num) + "<" + strconv.Itoa(s.depth) + ">"
}

func (s dummyState) GetPossibleActions() []mcts.Action {
	numPossibleActions := 5 - s.num
	if numPossibleActions < 0 {
		numPossibleActions = 0
	}
	possibleActions := make([]mcts.Action, numPossibleActions)
	for i := range possibleActions {
		possibleActions[i] = dummyAction{"act" + strconv.Itoa(i), i}
	}
	return possibleActions
}

func (s dummyState) GetSuccessorState(a mcts.Action) mcts.State {
	ac := a.(dummyAction)
	return dummyState{s.num + ac.i + 1, s.depth + 1}
}

func (s dummyState) EvaluateFinalState() float64 {
	return float64(s.num * s.depth)
}

// Dummy action
type dummyAction struct {
	s string
	i int
}

func (a dummyAction) String() string {
	return a.s + "<" + strconv.Itoa(a.i) + ">"
}

func main() {
	initState := dummyState{0, 0}
	explorationFactor := 10.0 //math.Sqrt(2)
	mcts := mcts.InitMCTS(initState, explorationFactor)
	for i := 0; i < 1000; i++ {
		mcts.RunIteration()
	}
	fmt.Println(mcts)
}
