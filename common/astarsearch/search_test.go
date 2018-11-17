package astarsearch

import (
	"testing"

	"github.com/RdecKa/bachleor-thesis/common/astarsearch"
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

func benchmarkAStarSearchOnHexGrid(actions []*hex.Action, size byte, b *testing.B) {
	state := hex.NewState(size, hex.Red)
	for _, a := range actions {
		s := state.GetSuccessorState(a).(hex.State)
		state = &s
	}

	for n := 0; n < b.N; n++ {
		initialState := hex.GetInitialState(state)
		aStarSearch := astarsearch.InitSearch(&initialState)
		aStarSearch.Search(false)
	}
}

/*
. . . . . . .
 . . . . . . .
  . . . . . . .
   . . . . . . .
    . . . . . . .
     . . . . . . .
      . . . . . . .
*/
func Benchmark0(b *testing.B) {
	actions := []*hex.Action{}

	benchmarkAStarSearchOnHexGrid(actions, 7, b)
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
func Benchmark1(b *testing.B) {
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

	benchmarkAStarSearchOnHexGrid(actions, 7, b)
}

/*
. . b . . r .
 . . b . r . .
  . . b b . . .
   r . b . . . .
    . b . b . r .
     . b . . r r .
      . b b r r r r
*/
func Benchmark2(b *testing.B) {
	actions := []*hex.Action{
		hex.NewAction(0, 3, hex.Red),
		hex.NewAction(2, 0, hex.Blue),
		hex.NewAction(5, 0, hex.Red),
		hex.NewAction(2, 1, hex.Blue),
		hex.NewAction(4, 1, hex.Red),
		hex.NewAction(2, 2, hex.Blue),
		hex.NewAction(3, 6, hex.Red),
		hex.NewAction(2, 3, hex.Blue),
		hex.NewAction(4, 6, hex.Red),
		hex.NewAction(1, 4, hex.Blue),
		hex.NewAction(5, 6, hex.Red),
		hex.NewAction(1, 5, hex.Blue),
		hex.NewAction(6, 6, hex.Red),
		hex.NewAction(1, 6, hex.Blue),
		hex.NewAction(4, 5, hex.Red),
		hex.NewAction(2, 6, hex.Blue),
		hex.NewAction(5, 4, hex.Red),
		hex.NewAction(3, 4, hex.Blue),
		hex.NewAction(5, 5, hex.Red),
		hex.NewAction(3, 2, hex.Blue),
	}

	benchmarkAStarSearchOnHexGrid(actions, 7, b)
}

/*
. . b b b r .
 . r b b r . .
  . . b b r . .
   r . b . . b .
    . b . b . r r
     . b . . r r r
      . b b r r r r
*/
func Benchmark3(b *testing.B) {
	actions := []*hex.Action{
		hex.NewAction(0, 3, hex.Red),
		hex.NewAction(2, 0, hex.Blue),
		hex.NewAction(5, 0, hex.Red),
		hex.NewAction(2, 1, hex.Blue),
		hex.NewAction(4, 1, hex.Red),
		hex.NewAction(2, 2, hex.Blue),
		hex.NewAction(3, 6, hex.Red),
		hex.NewAction(2, 3, hex.Blue),
		hex.NewAction(4, 6, hex.Red),
		hex.NewAction(1, 4, hex.Blue),
		hex.NewAction(5, 6, hex.Red),
		hex.NewAction(1, 5, hex.Blue),
		hex.NewAction(6, 6, hex.Red),
		hex.NewAction(1, 6, hex.Blue),
		hex.NewAction(4, 5, hex.Red),
		hex.NewAction(2, 6, hex.Blue),
		hex.NewAction(5, 4, hex.Red),
		hex.NewAction(3, 4, hex.Blue),
		hex.NewAction(5, 5, hex.Red),
		hex.NewAction(5, 3, hex.Blue),
		hex.NewAction(1, 1, hex.Red),
		hex.NewAction(3, 0, hex.Blue),
		hex.NewAction(4, 2, hex.Red),
		hex.NewAction(4, 0, hex.Blue),
		hex.NewAction(6, 4, hex.Red),
		hex.NewAction(3, 1, hex.Blue),
		hex.NewAction(6, 5, hex.Red),
		hex.NewAction(3, 2, hex.Blue),
	}

	benchmarkAStarSearchOnHexGrid(actions, 7, b)
}
