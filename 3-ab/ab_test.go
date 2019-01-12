package ab

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/RdecKa/0xAI/common/game/hex"
)

const patFileName = "../../common/game/hex/patterns.txt"

func benchmarkAB(actions []*hex.Action, size byte, b *testing.B) {
	state := hex.NewState(size, hex.Red)
	for _, a := range actions {
		s := state.GetSuccessorState(a).(hex.State)
		state = &s
	}

	gridChan, patChan, stopChan, resultChan := hex.CreatePatChecker(patFileName)
	defer func() { stopChan <- struct{}{} }()

	for n := 0; n < b.N; n++ {
		// Now when time is added, results cannot really be compared anymore ...
		AlphaBeta(state, time.Second, false, gridChan, patChan, resultChan,
			GetEstimateFunction("abLR"), "abLR")
	}
}

func _Benchmark0(b *testing.B) {
	actions := []*hex.Action{}

	benchmarkAB(actions, 7, b)
}

func _Benchmark1(b *testing.B) {
	actions := []*hex.Action{
		hex.NewAction(2, 2, hex.Red),
		hex.NewAction(3, 5, hex.Blue),
		hex.NewAction(1, 4, hex.Red),
		hex.NewAction(5, 4, hex.Blue),
	}

	benchmarkAB(actions, 7, b)
}

func _Benchmark2(b *testing.B) {
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

/*
. . . . . r r
 . . . b . . .
  . . . . . . .
   . . b . r r .
    . . . b . . .
     . b . . . . .
      . . . . . . .
*/
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

func benchAB(b *testing.B, depthLimit int, abSubtype string) {
	gridChan, patChan, stopChan, resultChan := hex.CreatePatChecker(patFileName)
	defer func() { stopChan <- struct{}{} }()

	_, state := getActionsAndStateSample()

	for n := 0; n < b.N; n++ {
		var oldTranspositionTable map[uint64]float64
		for depth := 2; depth <= depthLimit; depth += 2 {
			transpositionTable := make(map[uint64]float64)
			alphaBeta(context.TODO(), 0, depth, state, nil, math.Inf(-1), math.Inf(1),
				gridChan, patChan, resultChan, transpositionTable, oldTranspositionTable,
				false, GetEstimateFunction(abSubtype), abSubtype)
			oldTranspositionTable = transpositionTable
		}
	}
}

func BenchmarkAbLrLevel2(b *testing.B) {
	benchAB(b, 2, "abLR")
}

func BenchmarkAbLrLevel4(b *testing.B) {
	benchAB(b, 4, "abLR")
}

func BenchmarkAbLrLevel6(b *testing.B) {
	benchAB(b, 6, "abLR")
}

func BenchmarkAbDtLevel2(b *testing.B) {
	benchAB(b, 2, "abDT")
}

func BenchmarkAbDtLevel4(b *testing.B) {
	benchAB(b, 4, "abDT")
}

func BenchmarkAbDtLevel6(b *testing.B) {
	benchAB(b, 6, "abDT")
}
