package ab

import (
	"sort"

	"github.com/RdecKa/bachleor-thesis/common/game"
	"github.com/RdecKa/bachleor-thesis/common/game/hex"
)

// sortData is a collection used for sorting possibleActions
type sortData struct {
	data       []game.Action // slice of possible actions to be sorted
	dataValues []float64     // values of actions
	increasing bool          // true if the slice is to be sorted in increasing order, false for decreasing order
}

func initSortData(possibleActions []game.Action, state *hex.State,
	oldTransitionTable map[string]float64, increasing bool) *sortData {
	sd := &sortData{
		data:       possibleActions,
		dataValues: make([]float64, len(possibleActions)),
		increasing: increasing,
	}

	for i, a := range possibleActions {
		successorState := state.GetSuccessorState(a).(hex.State)
		tt, ok := oldTransitionTable[successorState.GetMapKey()]
		if !ok {
			sd.dataValues[i] = 0
		} else {
			sd.dataValues[i] = tt
		}
	}
	return sd
}

func (d *sortData) Len() int {
	return len(d.data)
}

func (d *sortData) Less(i, j int) bool {
	if d.increasing {
		return d.dataValues[i] < d.dataValues[j]
	}
	return d.dataValues[i] >= d.dataValues[j]
}

func (d *sortData) Swap(i, j int) {
	d.data[i], d.data[j] = d.data[j], d.data[i]
	d.dataValues[i], d.dataValues[j] = d.dataValues[j], d.dataValues[i]
}

func orderMoves(state *hex.State, possibleActions []game.Action,
	oldTransitionTable map[string]float64, increasing bool) []game.Action {

	if oldTransitionTable == nil {
		return possibleActions
	}

	sd := initSortData(possibleActions, state, oldTransitionTable, increasing)
	sort.Sort(sd)

	return sd.data
}
