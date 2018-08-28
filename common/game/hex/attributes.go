package hex

import (
	"fmt"

	"github.com/RdecKa/bachleor-thesis/common/game"
)

// All the attributes avaliable
var (
	AttrNumStones     = AttrNumberStones{}
	AttrOccRedRows    = AttrOccupiedRowsCols{Red, true}
	AttrOccRedCols    = AttrOccupiedRowsCols{Red, false}
	AttrOccBlueRows   = AttrOccupiedRowsCols{Blue, true}
	AttrOccBlueCols   = AttrOccupiedRowsCols{Blue, false}
	AttrPatCountRed0  = AttrPatternCount{Red, 0}
	AttrPatCountRed1  = AttrPatternCount{Red, 1}
	AttrPatCountRed2  = AttrPatternCount{Red, 2}
	AttrPatCountBlue0 = AttrPatternCount{Blue, 0}
	AttrPatCountBlue1 = AttrPatternCount{Blue, 1}
	AttrPatCountBlue2 = AttrPatternCount{Blue, 2}
)

// GenSamAttributes contains the attributes that are included in the sample
// generation
var GenSamAttributes = []game.Attribute{
	AttrNumStones,
	AttrOccRedRows,
	AttrOccRedCols,
	AttrOccBlueRows,
	AttrOccBlueCols,
	AttrPatCountRed0,
	AttrPatCountRed1,
	AttrPatCountRed2,
	AttrPatCountBlue0,
	AttrPatCountBlue1,
	AttrPatCountBlue2,
}

// ----------------------------
// |     AttrNumberStones     |
// ----------------------------

// AttrNumberStones takes care of getting the number of stones on the board
type AttrNumberStones struct{}

// GetAttributeName returns the name of an attribute
func (a AttrNumberStones) GetAttributeName() string {
	return "num_stones"
}

// GetAttributeValue returns the value of an attribute
func (a AttrNumberStones) GetAttributeValue(args []interface{}) int {
	s := args[0].(State)
	r, b, _ := s.GetNumOfStones()
	return r + b
}

// --------------------------------
// |     AttrOccupiedRowsCols     |
// --------------------------------

// AttrOccupiedRowsCols takes care of getting the number of occupied rows or
// columns for a player
type AttrOccupiedRowsCols struct {
	color Color // For which player rows/cols are counted
	rows  bool  // true: counting occupied rows; false: counting occupied columns
}

// GetAttributeName returns the name of an attribute
func (a AttrOccupiedRowsCols) GetAttributeName() string {
	n := "occ_"

	switch a.color {
	case Red:
		n += "red"
	case Blue:
		n += "blue"
	default:
		panic(fmt.Errorf("Invalid color %v", a.color))
	}

	n += "_"

	if a.rows {
		n += "rows"
	} else {
		n += "cols"
	}

	return n
}

// GetAttributeValue returns the value of an attribute
func (a AttrOccupiedRowsCols) GetAttributeValue(args []interface{}) int {
	patCount := args[1].([2][]int)
	i := -1
	switch a.color {
	case Red:
		i = 0
	case Blue:
		i = 1
	default:
		panic(fmt.Errorf("Invalid color %v", a.color))
	}

	var r int
	if a.rows {
		r = patCount[i][len(patCount[i])-2]
	} else {
		r = patCount[i][len(patCount[i])-1]
	}

	return r
}

// ----------------------------
// |     AttrPatternCount     |
// ----------------------------

// AttrPatternCount takes care of the number of a single pattern on the board
// for one player
type AttrPatternCount struct {
	color        Color // For which player patterns are counted
	patternIndex int   // Index of the pattern
}

// GetAttributeName returns the name of an attribute
func (a AttrPatternCount) GetAttributeName() string {
	var n string

	switch a.color {
	case Red:
		n = "red_p"
	case Blue:
		n = "blue_p"
	default:
		panic(fmt.Errorf("Invalid color %v", a.color))
	}

	return fmt.Sprintf("%s%d", n, a.patternIndex)
}

// GetAttributeValue returns the value of an attribute
func (a AttrPatternCount) GetAttributeValue(args []interface{}) int {
	patCount := args[1].([2][]int)
	i := -1
	switch a.color {
	case Red:
		i = 0
	case Blue:
		i = 1
	default:
		panic(fmt.Errorf("Invalid color %v", a.color))
	}

	return patCount[i][a.patternIndex]
}