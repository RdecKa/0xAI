package ab

import (
	"testing"
	"time"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

func benchmarkAB(actions []*hex.Action, size byte, b *testing.B) {
	state := hex.NewState(7, hex.Red)
	for _, a := range actions {
		s := state.GetSuccessorState(a).(hex.State)
		state = &s
	}

	for n := 0; n < b.N; n++ {
		// Now when time is added, results cannot really be compared anymore ...
		AlphaBeta(state, time.Second, "../../common/game/hex/patterns.txt")
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
