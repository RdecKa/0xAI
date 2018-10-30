package mcts

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/RdecKa/bachleor-thesis/common/game/hex"
	"github.com/RdecKa/bachleor-thesis/common/tree"
)

// ------------------
// |     record     |
// ------------------

type record struct {
	numChecked int          // number of already checked boards with the corresponding number of stones
	candidates []*tree.Node // candidates for continuing MCTS
}

func (r *record) addCandidate(c *tree.Node) {
	r.candidates = append(r.candidates, c)
}

func (r *record) oneChecked() {
	r.numChecked++
}

// -------------------------
// |     CandidateList     |
// -------------------------

// CandidateList stores the list of candidates for continuing MCTS
type CandidateList struct {
	list       []*record // Index in the array represents the number of stones on the board
	sortedList []*record // List of pointers to the same structs as list, but this list is ordered by record.numChecked
}

func (cl *CandidateList) String() string {
	stones := "# Stones    "
	checked := "# Checked   "
	cands := "# Candidates"
	for i, c := range cl.list {
		stones += fmt.Sprintf(" %5d", i)
		checked += fmt.Sprintf(" %5d", c.numChecked)
		cands += fmt.Sprintf(" %5d", len(c.candidates))
	}
	return fmt.Sprintf("%s\n%s\n%s\n", stones, checked, cands)
}

// NewCandidateList creates a new list of candidates.
func NewCandidateList(boardSize int) *CandidateList {
	list := make([]*record, boardSize*boardSize+1)
	for i := range list {
		list[i] = &record{
			numChecked: 0,
			candidates: make([]*tree.Node, 0),
		}
	}
	sortedList := make([]*record, boardSize*boardSize+1)
	for i := range sortedList {
		sortedList[i] = list[i]
	}
	return &CandidateList{
		list:       list,
		sortedList: sortedList,
	}
}

// AddCandidates accepts a list of new candidates and adds each of them in the
// suitable record (sublist) of CandidateList
func (cl *CandidateList) AddCandidates(newCandidates []*tree.Node) {
	for _, c := range newCandidates {
		state := c.GetValue().(*mctsNodeValue).GetState().(hex.State)
		r, b, _ := state.GetNumOfStones()
		ind := r + b
		cl.list[ind].addCandidate(c)
	}
}

// GetNextCandidateToExpand returns a node from which MCTS should be continued.
// The node is randomly chosen from the record that has been chosen the least
// number of times
//
// Possible improvement: Instead of calling sort.Sort, move the record towards
// the end of the array (after numChecked has been increased) until it is in a
// correct place again.
func (cl *CandidateList) GetNextCandidateToExpand() *tree.Node {
	sort.Sort(cl)
	i := 0
	for i < len(cl.sortedList) && len(cl.sortedList[i].candidates) == 0 {
		// If the sublist has no candidates, check the next sublist
		i++
	}
	if i == len(cl.sortedList) {
		return nil // No candidates in any of the sublists
	}
	selectedList := cl.sortedList[i].candidates
	nti := rand.Intn(len(selectedList))
	cl.sortedList[i].candidates = append(selectedList[:nti], selectedList[nti+1:]...)
	cl.sortedList[i].oneChecked()
	return selectedList[nti]
}

// IsEmpty returns true if all of the sublists (records) are empty, and false
// otherwise
func (cl *CandidateList) IsEmpty() bool {
	for _, c := range cl.list {
		if len(c.candidates) != 0 {
			return false
		}
	}
	return true
}

// Functions to implement sorting collection
func (cl *CandidateList) Len() int {
	return len(cl.sortedList)
}
func (cl *CandidateList) Less(i, j int) bool {
	return cl.sortedList[i].numChecked <= cl.sortedList[j].numChecked
}
func (cl *CandidateList) Swap(i, j int) {
	cl.sortedList[i], cl.sortedList[j] = cl.sortedList[j], cl.sortedList[i]
}
