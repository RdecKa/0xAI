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
	initState := hex.NewState(3)
	explorationFactor := math.Sqrt(2)
	mcts := mcts.InitMCTS(*initState, explorationFactor)
	for i := 0; i < 5; i++ {
		fmt.Println("Iteration")
		mcts.RunIteration()
		fmt.Println(mcts)
	}

	/*state := hex.NewState(5)

	actions := []*hex.Action{
		hex.NewAction(1, 0, hex.Red),
		hex.NewAction(0, 1, hex.Blue),
		hex.NewAction(1, 3, hex.Red),
		hex.NewAction(1, 1, hex.Blue),
		hex.NewAction(1, 2, hex.Red),
		hex.NewAction(2, 0, hex.Blue),
		hex.NewAction(2, 1, hex.Red),
		hex.NewAction(3, 0, hex.Blue),
		hex.NewAction(4, 0, hex.Red),
		hex.NewAction(3, 1, hex.Blue),
		hex.NewAction(4, 2, hex.Red),
		hex.NewAction(4, 1, hex.Blue),
		hex.NewAction(0, 0, hex.Red),
		hex.NewAction(0, 2, hex.Blue),
	}

	for _, a := range actions {
		state = state.GetSuccessorState(*a)
	}

	fmt.Println(state)

	fmt.Printf("Solution exists? %v\n", state.IsGoalState())*/
}
