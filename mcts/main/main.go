package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/RdecKa/mcts/game"
	"github.com/RdecKa/mcts/hex"
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

func (s dummyState) GetPossibleActions() []game.Action {
	numPossibleActions := 5 - s.num
	if numPossibleActions < 0 {
		numPossibleActions = 0
	}
	possibleActions := make([]game.Action, numPossibleActions)
	for i := range possibleActions {
		possibleActions[i] = dummyAction{"act" + strconv.Itoa(i), i}
	}
	return possibleActions
}

func (s dummyState) GetSuccessorState(a game.Action) game.State {
	ac := a.(dummyAction)
	return dummyState{s.num + ac.i + 1, s.depth + 1}
}

func (s dummyState) EvaluateGoalState() float64 {
	return float64(s.num * s.depth)
}

func (s dummyState) IsGoalState() bool {
	return false
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
	initState := hex.NewState(2)
	explorationFactor := math.Sqrt(2)
	mcts := mcts.InitMCTS(*initState, explorationFactor)
	for i := 0; i < 30000; i++ {
		mcts.RunIteration()
	}
	fmt.Println(mcts)
}
