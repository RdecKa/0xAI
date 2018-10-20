package ab

import (
	"context"
	"testing"
	"time"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

const patFileName = "../../common/game/hex/patterns.txt"

func benchmarkAB(actions []*hex.Action, size byte, b *testing.B) {
	state := hex.NewState(size, hex.Red)
	for _, a := range actions {
		s := state.GetSuccessorState(a).(hex.State)
		state = &s
	}

	gridChan, stopChan, resultChan := hex.CreatePatChecker(patFileName)
	defer func() { stopChan <- struct{}{} }()

	for n := 0; n < b.N; n++ {
		// Now when time is added, results cannot really be compared anymore ...
		AlphaBeta(state, time.Second, false, gridChan, resultChan)
	}
}

func Benchmark0(b *testing.B) {
	actions := []*hex.Action{}

	benchmarkAB(actions, 7, b)
}

func Benchmark1(b *testing.B) {
	actions := []*hex.Action{
		hex.NewAction(2, 2, hex.Red),
		hex.NewAction(3, 5, hex.Blue),
		hex.NewAction(1, 4, hex.Red),
		hex.NewAction(5, 4, hex.Blue),
	}

	benchmarkAB(actions, 7, b)
}

func Benchmark2(b *testing.B) {
	actions := []*hex.Action{
		hex.NewAction(5, 0, hex.Red),
		hex.NewAction(3, 1, hex.Blue),
		hex.NewAction(6, 0, hex.Red),
		hex.NewAction(2, 3, hex.Blue),
		hex.NewAction(4, 3, hex.Red),
		hex.NewAction(3, 4, hex.Blue),
		hex.NewAction(5, 3, hex.Red),
		hex.NewAction(1, 5, hex.Blue),
	}

	benchmarkAB(actions, 7, b)
}

func getActionsAndStateSample() ([]*hex.Action, *hex.State) {
	actions := []*hex.Action{
		hex.NewAction(5, 0, hex.Red),
		hex.NewAction(3, 1, hex.Blue),
		hex.NewAction(6, 0, hex.Red),
		hex.NewAction(2, 3, hex.Blue),
		hex.NewAction(4, 3, hex.Red),
		hex.NewAction(3, 4, hex.Blue),
		hex.NewAction(5, 3, hex.Red),
		hex.NewAction(1, 5, hex.Blue),
	}

	state := hex.NewState(7, hex.Red)
	for _, a := range actions {
		s := state.GetSuccessorState(a).(hex.State)
		state = &s
	}

	return actions, state
}

func benchAB(b *testing.B, depthLimit int) {
	gridChan, stopChan, resultChan := hex.CreatePatChecker(patFileName)
	defer func() { stopChan <- struct{}{} }()

	_, state := getActionsAndStateSample()

	for n := 0; n < b.N; n++ {
		var oldTranspositionTable map[string]float64
		for depth := 2; depth <= depthLimit; depth += 2 {
			transpositionTable := make(map[string]float64)
			alphaBeta(context.TODO(), 0, depth, state, nil, -1e14, 1e14, gridChan,
				resultChan, transpositionTable, oldTranspositionTable, false)
			oldTranspositionTable = transpositionTable
		}
	}
}

func BenchmarkABLevel2(b *testing.B) {
	benchAB(b, 2)
}

func BenchmarkABLevel4(b *testing.B) {
	benchAB(b, 4)
}

func BenchmarkABLevel6(b *testing.B) {
	benchAB(b, 6)
}
